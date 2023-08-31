package servers

import (
	"context"
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
}

type gateway struct {
	c *GatewayConfig
	h *http.Server

	placerCl api.Placer
	metaCl   api.MetaService
	dataCl   api.DataService
}

func NewGateway(c *GatewayConfig) api.App { return &gateway{c: c} }
func Gateway() api.App                    { return NewGateway(new(GatewayConfig)) }

func (*gateway) Name() string { return "gateway" }

func (g *gateway) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&g.c.Host, "host", DefaultGatewayHost, "listen host")
	f.StringVar(&g.c.Placer.RawPeers, "peers", DefaultPlacerRawPeers, "full placer peers")
}

func (g *gateway) Initialize() error {
	ps, err := peers.ParsePeers(g.c.Placer.RawPeers)
	if err != nil {
		return err
	}
	g.c.Placer.Peers = ps

	// TODO
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	g.h = &http.Server{Addr: g.c.Host, Handler: mux}

	return nil
}

func (g *gateway) Run() error                         { return g.h.ListenAndServe() }
func (g *gateway) Shutdown(ctx context.Context) error { return g.h.Shutdown(ctx) }
