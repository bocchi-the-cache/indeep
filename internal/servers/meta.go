package servers

import (
	"context"
	"flag"
	"net/http"

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
			Peers:         DefaultPlacerPeers,
			ClientTimeout: DefaultPeersIOTimeout,
		},
	}
}

type metaserver struct {
	config   *MetaserverConfig
	server   *http.Server
	placerCl api.Placer
}

func NewMetaserver(c *MetaserverConfig) api.Server { return &metaserver{config: c} }
func Metaserver() api.Server                       { return NewMetaserver(DefaultMetaserverConfig()) }

func (*metaserver) Name() string { return "metaserver" }

func (m *metaserver) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&m.config.Host, "host", api.DefaultMetaserverHost, "listen host")
	f.StringVar(&m.config.rawPlacerPeers, "placer-hosts", DefaultPlacerHosts.String(), "placer hosts URL")
}

func (m *metaserver) Setup() error {
	if m.config.rawPlacerPeers != "" {
		ps, err := peers.ParsePeers(m.config.rawPlacerPeers)
		if err != nil {
			return err
		}
		m.config.Placer.Peers = ps
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
