package logs

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/go-hclog"
)

const flags = log.LstdFlags | log.Lmicroseconds

func basic(prefix string) *log.Logger     { return log.New(out, prefix, flags) }
func shortFile(prefix string) *log.Logger { return log.New(out, prefix, flags|log.Lshortfile) }

var (
	out = os.Stderr

	slogCore   = basic("SLOG ")
	hcCore     = basic("HC ")
	badgerCore = shortFile("BADGER ")

	Std = shortFile("STD ")
	S   = Structured()
)

func Structured() *slog.Logger {
	return slog.New(slog.NewTextHandler(slogCore.Writer(), &slog.HandlerOptions{AddSource: true}))
}

func HashiCorp(name string) hclog.Logger {
	return hclog.FromStandardLogger(hcCore, &hclog.LoggerOptions{Name: name, IncludeLocation: true})
}

type badgerLog struct{ *log.Logger }

func Badger() badger.Logger { return &badgerLog{Logger: badgerCore} }

func (l *badgerLog) Errorf(f string, v ...interface{}) {
	_ = l.Output(3, fmt.Sprintf("[ERRO] "+f, v...))
}

func (l *badgerLog) Warningf(f string, v ...interface{}) {
	_ = l.Output(3, fmt.Sprintf("[WARN] "+f, v...))
}

func (l *badgerLog) Infof(f string, v ...interface{}) {
	_ = l.Output(3, fmt.Sprintf("[INFO] "+f, v...))
}

func (*badgerLog) Debugf(string, ...interface{}) {}
