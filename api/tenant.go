package api

import "errors"

type (
	AccessKey string
	SecretKey string
)

var ErrUnauthenticated = errors.New("unauthenticated")

type Tenants interface {
	SecretKey(ak AccessKey) (SecretKey, error)
	ListAll() ([]AccessKey, error)
}
