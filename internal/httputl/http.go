package httputl

import (
	"errors"
	"net/http"
	"strings"
)

const AuthorizationKey = "Authorization"

var (
	ErrEmptyAuthorization = errors.New("empty authorization")
	ErrEmptyCredential    = errors.New("empty credential")
)

type Authorization struct{ Scheme, Credential string }

func NewAuthorization(r *http.Request) (*Authorization, error) {
	raw := r.Header.Get(AuthorizationKey)
	if len(raw) == 0 {
		return nil, ErrEmptyAuthorization
	}

	raws := strings.SplitN(raw, " ", 2)
	if len(raws) != 2 {
		return nil, ErrEmptyCredential
	}

	return &Authorization{Scheme: raws[0], Credential: raws[1]}, nil
}
