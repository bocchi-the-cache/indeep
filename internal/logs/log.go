package logs

import (
	"log"
	"log/slog"
	"os"

	"github.com/hashicorp/go-hclog"
)

const coreFlags = log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix

func newCore(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, coreFlags)
}

var (
	E = Logger("ERR ")
	S = SLogger()

	slogCore = newCore("SLOG ")
	hcCore   = newCore("HC ")
)

func Logger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, coreFlags|log.Lshortfile)
}

func SLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(slogCore.Writer(), &slog.HandlerOptions{AddSource: true}))
}

func HcLogger(name string) hclog.Logger {
	return hclog.FromStandardLogger(hcCore, &hclog.LoggerOptions{Name: name, IncludeLocation: true})
}
