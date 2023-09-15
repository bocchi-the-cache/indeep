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
	"github.com/bocchi-the-cache/indeep/internal/hyped"
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

const (
	DefaultPlacerDataDir        = "placer-data"
	DefaultPlacerSnapshotRetain = 10
	DefaultPlacerLogCacheCap    = 128
	DefaultPlacerPeersConnPool  = 10
	DefaultPlacerPeersIOTimeout = 15 * time.Second

	PlacerLogDBFile    = "placer.log.bolt"
	PlacerStableDBFile = "placer.stable.bolt"
)

var (
	ErrPlacerUnknownID = errors.New("unknown placer ID")

	DefaultPlacerPeerMap = api.NewAddressMap(api.RaftScheme).Join(api.DefaultPlacerID, api.DefaultPlacerPeer)
)

type PlacerConfig struct {
	Host           string
	ID             raft.ServerID
	PeerMap        *api.AddressMap
	DataDir        string
	SnapshotRetain int
	LogCacheCap    int
	PeersConnPool  int
	PeersIOTimeout time.Duration

	rawPeers string
}

func DefaultPlacerConfig() *PlacerConfig {
	return &PlacerConfig{
		Host:           api.DefaultPlacerHost,
		ID:             api.DefaultPlacerID,
		PeerMap:        DefaultPlacerPeerMap,
		DataDir:        DefaultPlacerDataDir,
		SnapshotRetain: DefaultPlacerSnapshotRetain,
		LogCacheCap:    DefaultPlacerLogCacheCap,
		PeersConnPool:  DefaultPlacerPeersConnPool,
		PeersIOTimeout: DefaultPlacerPeersIOTimeout,
	}
}

func (c *PlacerConfig) hcLogger(name string) hclog.Logger {
	return logs.HcLogger(fmt.Sprintf("%s-%s", c.ID, name))
}

type placerServer struct {
	config *PlacerConfig
	peers  api.Peers
	server *http.Server
	rn     *raft.Raft
}

func NewPlacer(c *PlacerConfig) api.Server { return &placerServer{config: c} }
func Placer() api.Server                   { return NewPlacer(DefaultPlacerConfig()) }

func (*placerServer) Name() string { return "placer" }

func (s *placerServer) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&s.config.Host, "host", api.DefaultPlacerHost, "listen host")
	f.StringVar((*string)(&s.config.ID), "id", api.DefaultPlacerID, "placer ID")
	f.StringVar(&s.config.rawPeers, "peers", DefaultPlacerPeerMap.String(), "placer peers URL")
	f.StringVar(&s.config.DataDir, "data-dir", DefaultPlacerDataDir, "data directory")
	f.IntVar(&s.config.SnapshotRetain, "snap-retain", DefaultPlacerSnapshotRetain, "Raft snapshots to retain")
	f.IntVar(&s.config.LogCacheCap, "logcache-cap", DefaultPlacerLogCacheCap, "Raft log cache capacity")
	f.IntVar(&s.config.PeersConnPool, "conn-pool", DefaultPlacerPeersConnPool, "peer connections to pool")
	f.DurationVar(&s.config.PeersIOTimeout, "io-timeout", DefaultPlacerPeersIOTimeout, "peer IO timeout")
}

func (s *placerServer) Setup() error {
	if s.config.rawPeers != "" {
		ps, err := api.ParseAddressMap(s.config.rawPeers)
		if err != nil {
			return err
		}
		s.config.PeerMap = ps
	}
	s.peers = peers.NewPeers(s.config.PeerMap)

	p := s.peers.Lookup(s.config.ID)
	if p == nil {
		return fmt.Errorf("%w: peers=%s, id=%s", ErrPlacerUnknownID, s.config.PeerMap, s.config.ID)
	}

	config := raft.DefaultConfig()
	config.LocalID = s.config.ID
	config.Logger = s.config.hcLogger("raft")

	snaps, err := raft.NewFileSnapshotStoreWithLogger(
		s.config.DataDir,
		s.config.SnapshotRetain,
		s.config.hcLogger("snaps"),
	)
	if err != nil {
		return err
	}

	trans, err := raft.NewTCPTransportWithLogger(
		p.Address().Host,
		nil,
		s.config.PeersConnPool,
		s.config.PeersIOTimeout,
		s.config.hcLogger("trans"),
	)
	if err != nil {
		return err
	}

	logDB, err := raftboltdb.New(raftboltdb.Options{Path: filepath.Join(s.config.DataDir, PlacerLogDBFile)})
	if err != nil {
		return err
	}
	cachedLogDB, err := raft.NewLogCache(s.config.LogCacheCap, logDB)
	if err != nil {
		return err
	}

	stableDB, err := raftboltdb.New(raftboltdb.Options{Path: filepath.Join(s.config.DataDir, PlacerStableDBFile)})
	if err != nil {
		return err
	}

	rn, err := raft.NewRaft(config, s, cachedLogDB, stableDB, snaps, trans)
	if err != nil {
		return err
	}
	s.rn = rn
	s.rn.BootstrapCluster(s.peers.Configuration())

	s.server = &http.Server{
		Addr: s.config.Host,
		Handler: peers.
			ServeMux(p).
			HandleFunc(api.RpcGetMembers, hyped.Provider(s.HandleGetMembers)).
			HandleFunc(api.RpcAskLeader, hyped.Provider(s.HandleAskLeader)).
			Build(),
		ErrorLog: logs.E,
	}

	return nil
}

func (s *placerServer) ListenAndServe() error { return s.server.ListenAndServe() }
func (s *placerServer) Shutdown(ctx context.Context) error {
	return errors.Join(s.rn.Shutdown().Error(), s.server.Shutdown(ctx))
}

func (s *placerServer) HandleGetMembers() (raft.Configuration, error) {
	return s.GetMembers().Configuration(), nil
}

func (s *placerServer) HandleAskLeader() (api.Peer, error) {
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

func (s *placerServer) GetMembers() api.Peers { return s.peers }

func (s *placerServer) AskLeader(api.Peer) (api.Peer, error) {
	_, id := s.rn.LeaderWithID()
	return s.peers.Lookup(id), nil
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
