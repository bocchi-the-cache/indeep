package gateways

import "net/http"

func (g *gateway) Host() string { return g.config.Host }

func (g *gateway) DefineMux(mux *http.ServeMux) {
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
}

func (g *gateway) Close() error { return nil }
