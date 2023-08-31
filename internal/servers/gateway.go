package servers

import (
	"flag"
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

const DefaultGatewayHost = "127.0.0.1:11401"

type GatewayConfig struct {
	Host   string
	Placer clients.PlacerConfig

	rawPlacerPeers string
}

type gateway struct {
	config *GatewayConfig
	server *http.Server

	placerCl api.Placer
	metaCl   api.MetaService
	dataCl   api.DataService
}

func NewGateway(c *GatewayConfig) api.Server { return &gateway{config: c} }
func Gateway() api.Server                    { return NewGateway(new(GatewayConfig)) }

func (*gateway) Name() string { return "gateway" }

func (g *gateway) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&g.config.Host, "host", DefaultGatewayHost, "listen host")
	f.StringVar(&g.config.rawPlacerPeers, "peers", DefaultPlacerRawPeers, "full placer peers")
}

func (g *gateway) Setup() error {
	if g.config.Placer.Peers == nil {
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
	g.server = &http.Server{Addr: g.config.Host, Handler: mux}

	return nil
}

func (g *gateway) Server() *http.Server { return g.server }
