package servers

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

const (
	DefaultPlacerHost     = "127.0.0.1:11451"
	DefaultPlacerID       = "placer0"
	DefaultPlacerURL      = "tcp://127.0.0.1:11551"
	DefaultPlacerPeersURL = DefaultPlacerURL + peers.IDsPrefix + DefaultPlacerID
	DefaultSnapshotDir    = "."
	DefaultSnapshotRetain = 10
	DefaultLogDBFile      = "placer.log.bolt"
	DefaultLogCacheCap    = 128
	DefaultStableDBFile   = "placer.stable.bolt"
	DefaultPeersConnPool  = 10
	DefaultPeersIOTimeout = 15 * time.Second
)

var (
	ErrPlacerUnknownID = errors.New("unknown placer ID")

	DefaultPlacerPeers, _ = peers.ParsePeerURLs([2]string{DefaultPlacerID, DefaultPlacerURL})
	DefaultLogDBPath      = filepath.Join(DefaultSnapshotDir, DefaultLogDBFile)
	DefaultStableDBPath   = filepath.Join(DefaultSnapshotDir, DefaultStableDBFile)
)

type PlacerConfig struct {
	Host           string
	ID             raft.ServerID
	Peers          api.Peers
	SnapshotDir    string
	SnapshotRetain int
	LogDBPath      string
	LogCacheCap    int
	StableDBPath   string
	PeersConnPool  int
	PeersIOTimeout time.Duration

	rawPeers string
}

func DefaultPlacerConfig() *PlacerConfig {
	return &PlacerConfig{
		Host:           DefaultPlacerHost,
		ID:             DefaultPlacerID,
		Peers:          DefaultPlacerPeers,
		SnapshotDir:    DefaultSnapshotDir,
		SnapshotRetain: DefaultSnapshotRetain,
		LogDBPath:      DefaultLogDBPath,
		LogCacheCap:    DefaultLogCacheCap,
		StableDBPath:   DefaultStableDBPath,
		PeersConnPool:  DefaultPeersConnPool,
		PeersIOTimeout: DefaultPeersIOTimeout,
	}
}

func (c *PlacerConfig) hcLogger(name string) hclog.Logger {
	return logs.HcLogger(fmt.Sprintf("%s-%s", c.ID, name))
}

type placerServer struct {
	config *PlacerConfig
	server *http.Server
	fsm    *placerFSM
	rn     *raft.Raft
}

func NewPlacer(c *PlacerConfig) api.Server { return &placerServer{config: c} }
func Placer() api.Server                   { return NewPlacer(DefaultPlacerConfig()) }

func (*placerServer) Name() string { return "placer" }

func (s *placerServer) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&s.config.Host, "host", DefaultPlacerHost, "listen host")
	f.StringVar((*string)(&s.config.ID), "id", DefaultPlacerID, "placer ID")
	f.StringVar(&s.config.rawPeers, "peers", DefaultPlacerPeersURL, "placer peers URL")
	f.StringVar(&s.config.SnapshotDir, "snap-dir", DefaultSnapshotDir, "Raft snapshot base directory")
	f.IntVar(&s.config.SnapshotRetain, "snap-retain", DefaultSnapshotRetain, "Raft snapshots to retain")
	f.StringVar(&s.config.LogDBPath, "logdb", DefaultLogDBPath, "Raft log store path")
	f.IntVar(&s.config.LogCacheCap, "logcache-cap", DefaultLogCacheCap, "Raft log cache capacity")
	f.StringVar(&s.config.StableDBPath, "stabledb", DefaultStableDBPath, "Raft stable store path")
	f.IntVar(&s.config.PeersConnPool, "conn-pool", DefaultPeersConnPool, "peer connections to pool")
}

func (s *placerServer) Setup() error {
	if s.config.rawPeers != "" {
		ps, err := peers.ParsePeers(s.config.rawPeers)
		if err != nil {
			return err
		}
		s.config.Peers = ps
	}

	p := s.config.Peers.Lookup(s.config.ID)
	if p == nil {
		return fmt.Errorf("%w: peers=%s, id=%s", ErrPlacerUnknownID, s.config.Peers, s.config.ID)
	}

	config := raft.DefaultConfig()
	config.LocalID = s.config.ID
	config.Logger = s.config.hcLogger("raft")

	snaps, err := raft.NewFileSnapshotStoreWithLogger(
		s.config.SnapshotDir,
		s.config.SnapshotRetain,
		s.config.hcLogger("snaps"),
	)
	if err != nil {
		return err
	}

	trans, err := raft.NewTCPTransportWithLogger(
		p.URL().Host,
		nil,
		s.config.PeersConnPool,
		s.config.PeersIOTimeout,
		s.config.hcLogger("trans"),
	)
	if err != nil {
		return err
	}

	// TODO
	s.fsm = &placerFSM{
		peers:    s.config.Peers,
		self:     p,
		leader:   p,
		isLeader: true,
	}

	logDB, err := raftboltdb.New(raftboltdb.Options{Path: s.config.LogDBPath})
	if err != nil {
		return err
	}
	cachedLogDB, err := raft.NewLogCache(s.config.LogCacheCap, logDB)
	if err != nil {
		return err
	}

	stableDB, err := raftboltdb.New(raftboltdb.Options{Path: s.config.StableDBPath})
	if err != nil {
		return err
	}

	rn, err := raft.NewRaft(config, s.fsm, cachedLogDB, stableDB, snaps, trans)
	if err != nil {
		return err
	}
	s.rn = rn
	s.rn.BootstrapCluster(s.config.Peers.Configuration())

	mux := http.NewServeMux()
	mux.HandleFunc(p.RPC(api.RpcGetMembers).Path, s.Members)
	mux.HandleFunc(p.RPC(api.RpcAskLeader).Path, s.Leader)
	mux.HandleFunc(p.RPC(api.RpcLookupMetaService).Path, s.LookupMetaService)
	mux.HandleFunc(p.RPC(api.RpcAddMetaService).Path, s.AddMetaService)
	mux.HandleFunc(p.RPC(api.RpcLookupDataService).Path, s.LookupDataService)
	mux.HandleFunc(p.RPC(api.RpcAddDataService).Path, s.AddDataService)
	s.server = &http.Server{
		Addr:     s.config.Host,
		Handler:  mux,
		ErrorLog: logs.E,
	}

	return nil
}

func (s *placerServer) ListenAndServe() error { return s.server.ListenAndServe() }
func (s *placerServer) Shutdown(ctx context.Context) error {
	return errors.Join(s.rn.Shutdown().Error(), s.server.Shutdown(ctx))
}

func (s *placerServer) Members(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	jsonhttp.W(w).OK(s.fsm.Members())
}

func (s *placerServer) Leader(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	jw := jsonhttp.W(w)
	l, err := s.fsm.Leader(nil)
	if err != nil {
		jw.Err(err)
		return
	}
	jw.OK(l)
}

func (s *placerServer) LookupMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = s.fsm.LookupMetaService(nil)
}

func (s *placerServer) AddMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = s.fsm.AddMetaService()
}

func (s *placerServer) LookupDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = s.fsm.LookupDataService(nil)
}

func (s *placerServer) AddDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = s.fsm.AddDataService()
}

type placerFSM struct {
	// FIXME: Any locks here?
	peers    api.Peers
	self     api.Peer
	leader   api.Peer
	isLeader bool
}

func (p *placerFSM) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) Members() api.Peers {
	return p.peers
}

func (p *placerFSM) Leader(api.Peer) (api.Peer, error) {
	return p.leader, nil
}

func (p *placerFSM) LookupMetaService(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) AddMetaService() error {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) LookupDataService(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) AddDataService() error {
	//TODO implement me
	panic("implement me")
}
