package metaservers

import (
	"flag"
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

type metaserver struct {
	config   *MetaserverConfig
	peers    api.Peers
	server   *http.Server
	placerCl api.Placer
	mux      api.StreamLayerMux
}

func NewMetaserver(c *MetaserverConfig) api.Server { return &metaserver{config: c} }
func Metaserver() api.Server                       { return NewMetaserver(DefaultMetaserverConfig()) }

func (*metaserver) Name() string { return "metaserver" }

func (m *metaserver) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&m.config.Host, "host", api.DefaultMetaserverHost, "listen host")
	f.StringVar((*string)(&m.config.ID), "id", api.DefaultMetaserverID, "placer ID")
	f.StringVar(&m.config.rawPeers, "peers", api.DefaultMetaserverMultipeerMap.String(), "metaserver peers URL")
	f.StringVar(&m.config.DataDir, "data-dir", DefaultMetaserverDataDir, "data directory")
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
	ps, err := peers.NewPeers(m.config.PeerMap)
	if err != nil {
		return err
	}
	m.peers = ps

	if _, err := m.peers.Lookup(m.config.ID); err != nil {
		return err
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

	return nil
}

func (m *metaserver) Host() string { return m.config.Host }

func (m *metaserver) Close() error { return nil }

func (m *metaserver) DefineMux(mux *http.ServeMux) {
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
}
