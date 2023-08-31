package servers

import "github.com/bocchi-the-cache/indeep/internal/clients"

type Config struct {
	Gateway GatewayConfig
}

type GatewayConfig struct {
	Addr string

	Placer clients.PlacerConfig
}
