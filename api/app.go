package api

import (
	"context"
	"flag"
)

type App interface {
	Name() string
	DefineFlags(f *flag.FlagSet)
	Setup() error
}

type Server interface {
	App

	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
