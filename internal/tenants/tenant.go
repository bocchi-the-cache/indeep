package tenants

import (
	"github.com/bocchi-the-cache/indeep/api"
)

const (
	DefaultDisplayName = "Admin"
	DefaultAccessKey   = "root"
	DefaultSecretKey   = "secret"
)

type tenants struct{}

func New() api.Tenants { return new(tenants) }

func (*tenants) Get(api.AccessKey) (*api.Tenant, error) {
	return &api.Tenant{
		DisplayName: DefaultDisplayName,
		AccessKey:   DefaultAccessKey,
		SecretKey:   DefaultSecretKey,
	}, nil
}

func (*tenants) ListAll() ([]api.AccessKey, error) { return []api.AccessKey{DefaultAccessKey}, nil }
