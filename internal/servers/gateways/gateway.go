package gateways

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

type GatewayConfig struct {
	Host   string
	Placer clients.PlacerConfig

	rawPlacerHosts string
}

func DefaultGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Host: api.DefaultGatewayHost,
		Placer: clients.PlacerConfig{
			HostMap:       api.DefaultPlacerHostMap,
			ClientTimeout: 15 * time.Second,
		},
	}
}

type gateway struct {
	config   *GatewayConfig
	peers    api.Peers
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
	f.StringVar(&g.config.rawPlacerHosts, "placer-hosts", api.DefaultPlacerHostMap.String(), "placer hosts URL")
}

func (g *gateway) Setup() error {
	if g.config.rawPlacerHosts != "" {
		ps, err := api.ParseAddressMap(g.config.rawPlacerHosts)
		if err != nil {
			return err
		}
		g.config.Placer.HostMap = ps
	}
	g.peers = peers.NewPeers(g.config.Placer.HostMap)

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
