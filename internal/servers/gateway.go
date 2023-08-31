package servers

import (
	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
)

type gateway struct {
	placerCl api.Placer
	metaCl   api.MetaService
	dataCl   api.DataService
}

func ServeGateway(c *Config) error {
	placerCl, err := clients.NewPlacer(&c.Gateway.Placer)
	if err != nil {
		return err
	}

	g := &gateway{
		placerCl: placerCl,
	}

	_ = g // TODO
	return nil
}
