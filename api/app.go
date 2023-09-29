package api

import (
	"flag"
	"io"
	"net/http"
)

type App interface {
	Name() string
	DefineFlags(f *flag.FlagSet)
	Setup() error
}

type Server interface {
	App

	io.Closer
	Host() string
	DefineMux(mux *http.ServeMux)
}
