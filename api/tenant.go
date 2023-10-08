package api

type (
	AccessKey string
	SecretKey string
)

type (
	Tenant struct {
		DisplayName string
		AccessKey   AccessKey
		SecretKey   SecretKey
	}

	Tenants interface {
		Get(ak AccessKey) (*Tenant, error)
		ListAll() ([]AccessKey, error)
	}
)
