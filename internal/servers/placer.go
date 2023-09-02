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
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
	"github.com/bocchi-the-cache/indeep/internal/typedh"
)

const (
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

	DefaultPlacerPeers  = peers.RaftPeers().Join(api.DefaultPlacerID, peers.Voter(api.DefaultPlacerPeer))
	DefaultPlacerHosts  = peers.HostPeers().Join(api.DefaultPlacerID, peers.Voter(api.DefaultPlacerHost))
	DefaultLogDBPath    = filepath.Join(DefaultSnapshotDir, DefaultLogDBFile)
	DefaultStableDBPath = filepath.Join(DefaultSnapshotDir, DefaultStableDBFile)
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
		Host:           api.DefaultPlacerHost,
		ID:             api.DefaultPlacerID,
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
	rn     *raft.Raft
}

func NewPlacer(c *PlacerConfig) api.Server { return &placerServer{config: c} }
func Placer() api.Server                   { return NewPlacer(DefaultPlacerConfig()) }

func (*placerServer) Name() string { return "placer" }

func (s *placerServer) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&s.config.Host, "host", api.DefaultPlacerHost, "listen host")
	f.StringVar((*string)(&s.config.ID), "id", api.DefaultPlacerID, "placer ID")
	f.StringVar(&s.config.rawPeers, "peers", DefaultPlacerPeers.String(), "placer peers URL")
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

	rn, err := raft.NewRaft(config, s, cachedLogDB, stableDB, snaps, trans)
	if err != nil {
		return err
	}
	s.rn = rn
	s.rn.BootstrapCluster(s.config.Peers.Configuration())

	s.server = &http.Server{
		Addr: s.config.Host,
		Handler: peers.
			Mux(p).
			HandleFunc(api.RpcGetMembers, typedh.Provider(s.HandleGetMembers)).
			HandleFunc(api.RpcAskLeader, typedh.Provider(s.HandleAskLeader)).
			Build(),
		ErrorLog: logs.E,
	}

	return nil
}

func (s *placerServer) ListenAndServe() error { return s.server.ListenAndServe() }
func (s *placerServer) Shutdown(ctx context.Context) error {
	return errors.Join(s.rn.Shutdown().Error(), s.server.Shutdown(ctx))
}

func (s *placerServer) HandleGetMembers() (*raft.Configuration, error) {
	conf := s.GetMembers().Configuration()
	return &conf, nil
}

func (s *placerServer) HandleAskLeader() (*api.PeerInfo, error) {
	return s.AskLeader(nil)
}

func (s *placerServer) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) GetMembers() api.Peers { return s.config.Peers }

func (s *placerServer) AskLeader(api.Peer) (*api.PeerInfo, error) {
	leader, id := s.rn.LeaderWithID()
	return &api.PeerInfo{ID: id, Peer: peers.Voter(leader)}, nil
}

func (s *placerServer) LookupMetaService(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) AddMetaService() error {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) LookupDataService(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) AddDataService() error {
	//TODO implement me
	panic("implement me")
}
