package apps

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/logs"
)

func RunServer(a api.App, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := Run(a, args); !errors.Is(err, http.ErrServerClosed) {
			logs.E.Fatal(err)
		}
		cancel()
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	if err := a.Shutdown(context.Background()); err != nil {
		logs.E.Fatal(err)
	}
	<-ctx.Done()
}

func Run(a api.App, args []string) error {
	f := a.FlagSet()
	if err := f.Parse(args); err != nil {
		return err
	}
	if err := a.Initialize(); err != nil {
		return err
	}
	return a.Run()
}
