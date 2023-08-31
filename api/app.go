package api

import (
	"flag"
	"net/http"
)

type App interface {
	Name() string
	DefineFlags(f *flag.FlagSet)
	Setup() error
}

type Server interface {
	App
	Server() *http.Server
}
