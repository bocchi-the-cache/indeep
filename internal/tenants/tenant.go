package tenants

import (
	"github.com/bocchi-the-cache/indeep/api"
)

const (
	DefaultAccessKey = "root"
	DefaultSecretKey = "secret"
)

type tenants struct{}

func New() api.Tenants                                          { return new(tenants) }
func (*tenants) SecretKey(api.AccessKey) (api.SecretKey, error) { return DefaultSecretKey, nil }
func (*tenants) ListAll() ([]api.AccessKey, error)              { return []api.AccessKey{DefaultAccessKey}, nil }
