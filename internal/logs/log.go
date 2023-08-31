package logs

import (
	"log"
	"os"
)

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
}

var (
	I = newLogger("INF ")
	E = newLogger("ERR ")
)
