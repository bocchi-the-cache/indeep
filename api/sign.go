package api

import (
	"net/http"
)

type SigChecker interface {
	CheckSigV4(r *http.Request) (bool, error)
}
