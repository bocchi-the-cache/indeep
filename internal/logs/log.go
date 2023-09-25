package logs

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/hashicorp/go-hclog"
)

const coreFlags = log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix

func newCore(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, coreFlags)
}

var (
	E = Logger("ERR ")
	S = Structured()

	slogCore   = newCore("SLOG ")
	hcCore     = newCore("HC ")
	pebbleCore = newCore("PEBBLE ")
)

func Logger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, coreFlags|log.Lshortfile)
}

func Structured() *slog.Logger {
	return slog.New(slog.NewTextHandler(slogCore.Writer(), &slog.HandlerOptions{AddSource: true}))
}

func HashiCorp(name string) hclog.Logger {
	return hclog.FromStandardLogger(hcCore, &hclog.LoggerOptions{Name: name, IncludeLocation: true})
}

type pebbleLog struct{ *log.Logger }

func (p *pebbleLog) Infof(format string, args ...interface{}) {
	_ = p.Logger.Output(2, fmt.Sprintf(format, args...))
}

func (p *pebbleLog) Fatalf(format string, args ...interface{}) {
	_ = p.Logger.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Pebble() pebble.Logger { return &pebbleLog{pebbleCore} }
