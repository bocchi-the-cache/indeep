package api

import (
	"context"
	"flag"
)

type App interface {
	Name() string
	DefineFlags(f *flag.FlagSet)
	Initialize() error
	Run() error
	Shutdown(ctx context.Context) error
}
