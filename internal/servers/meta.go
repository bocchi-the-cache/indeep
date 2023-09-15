package servers

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

type MetaserverConfig struct {
	Host   string
	Placer clients.PlacerConfig

	rawPlacerPeers string
}

func DefaultMetaserverConfig() *MetaserverConfig {
	return &MetaserverConfig{
		Host: api.DefaultMetaserverHost,
		Placer: clients.PlacerConfig{
			PeerMap:       DefaultPlacerPeerMap,
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
	f.StringVar(&m.config.rawPlacerPeers, "placer-hosts", DefaultPlacerHostMap.String(), "placer hosts URL")
}

func (m *metaserver) Setup() error {
	if m.config.rawPlacerPeers != "" {
		ps, err := api.ParseAddressMap(m.config.rawPlacerPeers)
		if err != nil {
			return err
		}
		m.config.Placer.PeerMap = ps
	}
	m.peers = peers.NewPeers(m.config.Placer.PeerMap)

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
