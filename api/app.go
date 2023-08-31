package api

import (
	"context"
	"flag"
)

type App interface {
	FlagSet() *flag.FlagSet
	Initialize() error
	Run() error
	Shutdown(ctx context.Context) error
}
