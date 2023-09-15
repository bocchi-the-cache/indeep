package gateways

import "context"

func (g *gateway) ListenAndServe() error              { return g.server.ListenAndServe() }
func (g *gateway) Shutdown(ctx context.Context) error { return g.server.Shutdown(ctx) }
