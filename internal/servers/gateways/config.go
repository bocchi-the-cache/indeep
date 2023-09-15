package gateways

import (
	"time"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
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
