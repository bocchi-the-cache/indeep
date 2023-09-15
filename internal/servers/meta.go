package servers

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

const (
	DefaultMetaserverDataDir        = "metaserver-data"
	DefaultMetaserverSnapshotRetain = 10
	DefaultMetaserverLogCacheCap    = 128
	DefaultMetaserverPeersConnPool  = 10
	DefaultMetaserverPeersIOTimeout = 15 * time.Second
)

var (
	ErrMetaserverUnknownID = errors.New("unknown metaserver ID")

	DefaultMetaserverMultipeerMap = api.NewAddressMap(api.RaftScheme).Join(api.DefaultMetaserverID, api.DefaultMetaServerMultiPeer)
)

type MetaserverConfig struct {
	Host           string
	ID             raft.ServerID
	PeerMap        *api.AddressMap
	DataDir        string
	SnapshotRetain int
	LogCacheCap    int
	PeersConnPool  int
	PeersIOTimeout time.Duration
	Placer         clients.PlacerConfig

	rawPeers       string
	rawPlacerHosts string
}

func DefaultMetaserverConfig() *MetaserverConfig {
	return &MetaserverConfig{
		Host: api.DefaultMetaserverHost,
		Placer: clients.PlacerConfig{
			HostMap:       api.DefaultPlacerHostMap,
			ClientTimeout: 15 * time.Second,
		},
	}
}

type metaserver struct {
	config *MetaserverConfig
	peers  api.Peers
	server *http.Server

	placerCl api.Placer

	mux api.StreamLayerMux
}

func NewMetaserver(c *MetaserverConfig) api.Server { return &metaserver{config: c} }
func Metaserver() api.Server                       { return NewMetaserver(DefaultMetaserverConfig()) }

func (*metaserver) Name() string { return "metaserver" }

func (m *metaserver) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&m.config.Host, "host", api.DefaultMetaserverHost, "listen host")
	f.StringVar((*string)(&m.config.ID), "id", api.DefaultMetaserverID, "placer ID")
	f.StringVar(&m.config.rawPeers, "peers", DefaultMetaserverMultipeerMap.String(), "metaserver peers URL")
	f.StringVar(&m.config.DataDir, "data-dir", DefaultMetaserverDataDir, "data directory")
	f.IntVar(&m.config.SnapshotRetain, "snap-retain", DefaultMetaserverSnapshotRetain, "Raft snapshots to retain")
	f.IntVar(&m.config.LogCacheCap, "logcache-cap", DefaultMetaserverLogCacheCap, "Raft log cache capacity")
	f.IntVar(&m.config.PeersConnPool, "conn-pool", DefaultMetaserverPeersConnPool, "peer connections to pool")
	f.DurationVar(&m.config.PeersIOTimeout, "io-timeout", DefaultMetaserverPeersIOTimeout, "peer IO timeout")
	f.StringVar(&m.config.rawPlacerHosts, "placer-hosts", api.DefaultPlacerHostMap.String(), "placer hosts URL")
}

func (m *metaserver) Setup() error {
	if m.config.rawPeers != "" {
		ps, err := api.ParseAddressMap(m.config.rawPeers)
		if err != nil {
			return err
		}
		m.config.PeerMap = ps
	}
	m.peers = peers.NewPeers(m.config.PeerMap)

	p := m.peers.Lookup(m.config.ID)
	if p == nil {
		return fmt.Errorf("%w: peers=%s, id=%s", ErrMetaserverUnknownID, m.config.PeerMap, m.config.ID)
	}

	if m.config.rawPlacerHosts != "" {
		ps, err := api.ParseAddressMap(m.config.rawPlacerHosts)
		if err != nil {
			return err
		}
		m.config.Placer.HostMap = ps
	}
	placerCl, err := clients.NewPlacer(&m.config.Placer)
	if err != nil {
		return err
	}
	m.placerCl = placerCl

	// TODO
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	m.server = &http.Server{
		Addr:     m.config.Host,
		Handler:  mux,
		ErrorLog: logs.E,
	}

	return nil
}

func (m *metaserver) ListenAndServe() error              { return m.server.ListenAndServe() }
func (m *metaserver) Shutdown(ctx context.Context) error { return m.server.Shutdown(ctx) }

func (m *metaserver) Lookup(key api.MetaKey) (api.MetaPartition, error) {
	//TODO implement me
	panic("implement me")
}

func (m *metaserver) logDBPath(groupID api.GroupID) string {
	return filepath.Join(m.config.DataDir, fmt.Sprintf("metaserver.log.%s.bolt", groupID))
}

func (m *metaserver) stableDBPath(groupID api.GroupID) string {
	return filepath.Join(m.config.DataDir, fmt.Sprintf("metaserver.stable.%s.bolt", groupID))
}
