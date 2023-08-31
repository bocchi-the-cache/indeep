package servers

import (
	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
)

type GatewayConfig struct {
	Addr string

	Placer clients.PlacerConfig
}

type gateway struct {
	placerCl api.Placer
	metaCl   api.MetaService
	dataCl   api.DataService
}

func ServeGateway(c *GatewayConfig) error {
	placerCl, err := clients.NewPlacer(&c.Placer)
	if err != nil {
		return err
	}

	g := &gateway{
		placerCl: placerCl,
	}

	_ = g // TODO
	return nil
}
