package logs

import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
)

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, log.LstdFlags|log.Lshortfile|log.Lmicroseconds|log.Lmsgprefix)
}

var (
	I = newLogger("INF ")
	E = newLogger("ERR ")

	lib = newLogger("LIB ")
)

func HcLogger(name string) hclog.Logger {
	return hclog.FromStandardLogger(lib, &hclog.LoggerOptions{Name: name})
}
