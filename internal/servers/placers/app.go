package placers

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/databases"
	"github.com/bocchi-the-cache/indeep/internal/hyped"
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

var ErrPlacerUnknownID = errors.New("unknown placer ID")

type placerServer struct {
	*databases.DB
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
	f.StringVar(&s.config.rawPeers, "peers", api.DefaultPlacerPeerMap.String(), "placer peers URL")
	f.StringVar(&s.config.DataDir, "data-dir", DefaultPlacerDataDir, "data directory")
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

	stableDB, err := raftboltdb.New(raftboltdb.Options{Path: filepath.Join(s.config.DataDir, PlacerStableDBFile)})
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

	db, err := databases.Open(filepath.Join(s.config.DataDir, "db"))
	if err != nil {
		return err
	}
	s.DB = db

	rn, err := raft.NewRaft(config, s, cachedLogDB, stableDB, s, trans)
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
