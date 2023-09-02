package apps

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/logs"
)

func Setup(a api.App, args []string) error {
	f := flag.NewFlagSet(a.Name(), flag.ContinueOnError)
	a.DefineFlags(f)
	if err := f.Parse(args); err != nil {
		return err
	}
	return a.Setup()
}

func MainServer(s api.Server) { RunServer(s, os.Args[1:]) }

func RunServer(s api.Server, args []string) {
	if err := Setup(s, args); err != nil {
		logs.S.Error("setup error", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logs.S.Error("listen error", "err", err)
		}
		cancel()
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	if err := s.Shutdown(context.Background()); err != nil {
		logs.S.Error("shutdown error", "err", err)
	}
	<-ctx.Done()
}
