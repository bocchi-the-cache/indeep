package tenants

import (
	"fmt"

	"github.com/bocchi-the-cache/indeep/api"
)

const (
	DefaultAccessKey = "root"
	DefaultSecretKey = "secret"
)

type tenants struct{}

func New() api.Tenants { return new(tenants) }

func (*tenants) Authenticate(ak api.AccessKey, sk api.SecretKey) error {
	if ak == DefaultAccessKey && sk == DefaultSecretKey {
		return nil
	}
	return fmt.Errorf("%w: ak=%s", api.ErrUnauthenticated, ak)
}

func (*tenants) ListAll() ([]api.AccessKey, error) { return []api.AccessKey{DefaultAccessKey}, nil }
