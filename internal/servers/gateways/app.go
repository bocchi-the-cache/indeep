package gateways

import (
	"flag"
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
	"github.com/bocchi-the-cache/indeep/internal/sigv4"
	"github.com/bocchi-the-cache/indeep/internal/tenants"
	"github.com/bocchi-the-cache/indeep/internal/utils/hyped"
)

type gateway struct {
	config   *GatewayConfig
	codec    hyped.Codec
	sigChk   api.SigChecker
	server   *http.Server
	placerCl api.Placer
	metaCl   api.MetaService
	dataCl   api.DataService
}

func NewGateway(c *GatewayConfig) api.Server {
	codec := hyped.XML()
	return &gateway{
		config: c,
		codec:  codec,
		sigChk: sigv4.New(tenants.New(), codec),
	}
}

func Gateway() api.Server { return NewGateway(DefaultGatewayConfig()) }

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
	placerCl, err := clients.NewPlacer(&g.config.Placer)
	if err != nil {
		return err
	}
	g.placerCl = placerCl

	return nil
}

func (g *gateway) Host() string { return g.config.Host }

func (g *gateway) DefineMux(mux *http.ServeMux) {
	mux.HandleFunc("/", g.sigChk.WithSigV4(hyped.ProviderWith(g.codec, g.ListBuckets)))
}

func (g *gateway) Close() error { return nil }
