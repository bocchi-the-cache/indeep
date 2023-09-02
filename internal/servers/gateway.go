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

type GatewayConfig struct {
	Host   string
	Placer clients.PlacerConfig

	rawPlacerPeers string
}

func DefaultGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Host: api.DefaultGatewayHost,
		Placer: clients.PlacerConfig{
			Peers:         DefaultPlacerPeers,
			ClientTimeout: DefaultPeersIOTimeout,
		},
	}
}

type gateway struct {
	config   *GatewayConfig
	server   *http.Server
	placerCl api.Placer
	metaCl   api.MetaService
	dataCl   api.DataService
}

func NewGateway(c *GatewayConfig) api.Server { return &gateway{config: c} }
func Gateway() api.Server                    { return NewGateway(DefaultGatewayConfig()) }

func (*gateway) Name() string { return "gateway" }

func (g *gateway) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&g.config.Host, "host", api.DefaultGatewayHost, "listen host")
	f.StringVar(&g.config.rawPlacerPeers, "placer-hosts", DefaultPlacerHosts.String(), "placer hosts URL")
}

func (g *gateway) Setup() error {
	if g.config.rawPlacerPeers != "" {
		ps, err := peers.ParsePeers(g.config.rawPlacerPeers)
		if err != nil {
			return err
		}
		g.config.Placer.Peers = ps
	}

	placerCl, err := clients.NewPlacer(&g.config.Placer)
	if err != nil {
		return err
	}
	g.placerCl = placerCl

	// TODO
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	g.server = &http.Server{
		Addr:     g.config.Host,
		Handler:  mux,
		ErrorLog: logs.E,
	}

	return nil
}

func (g *gateway) ListenAndServe() error              { return g.server.ListenAndServe() }
func (g *gateway) Shutdown(ctx context.Context) error { return g.server.Shutdown(ctx) }
