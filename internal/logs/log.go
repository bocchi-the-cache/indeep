package logs

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dgraph-io/badger/v4"

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
	badgerCore = newCore("BADGER ")
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

type badgerLog struct{ *log.Logger }

func Badger() badger.Logger { return &badgerLog{Logger: badgerCore} }

func (l *badgerLog) Errorf(f string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf("[ERRO] "+f, v...))
}

func (l *badgerLog) Warningf(f string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf("[WARN] "+f, v...))
}

func (l *badgerLog) Infof(f string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf("[INFO] "+f, v...))
}

func (*badgerLog) Debugf(string, ...interface{}) {}
