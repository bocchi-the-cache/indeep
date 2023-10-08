package gateways

import "github.com/bocchi-the-cache/indeep/api"

var _ = (api.Gateway)((*gateway)(nil))

func (g *gateway) ListBuckets() ([]string, error) {
	// TODO
	return nil, nil
}
