package api

import (
	"errors"
	"net/http"
)

var ErrUnknownAuthScheme = errors.New("unknown authorization scheme")

type SigV4Checker interface {
	CheckSigV4(r *http.Request) error
}
