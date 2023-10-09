package httputl

import (
	"errors"
	"net/http"
	"strings"
)

const AuthorizationKey = "Authorization"

var (
	ErrEmptyAuthorization = errors.New("empty header authorization")
	ErrEmptyCredential    = errors.New("empty header credential")
)

type Authorization struct{ Scheme, Credential string }

func NewAuthorization(header http.Header) (*Authorization, error) {
	raw := header.Get(AuthorizationKey)
	if len(raw) == 0 {
		return nil, ErrEmptyAuthorization
	}

	raws := strings.SplitN(raw, " ", 2)
	if len(raws) != 2 {
		return nil, ErrEmptyCredential
	}

	return &Authorization{Scheme: raws[0], Credential: raws[1]}, nil
}

type Router interface {
	Route(r *http.Request) http.HandlerFunc
}

func NotImplemented(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
