package gateways

import "github.com/bocchi-the-cache/indeep/api"

var _ = (api.Gateway)((*gateway)(nil))

func (g *gateway) ListBuckets() (*api.ListAllMyBucketsResult, error) {
	// TODO
	return new(api.ListAllMyBucketsResult), nil
}
