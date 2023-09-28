package api

import "errors"

type (
	AccessKey string
	SecretKey string
)

var ErrUnauthenticated = errors.New("unauthenticated")

type Tenants interface {
	Authenticate(ak AccessKey, sk SecretKey) error
	ListAll() ([]AccessKey, error)
}
