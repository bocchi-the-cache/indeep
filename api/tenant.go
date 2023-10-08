package api

type (
	AccessKey string
	SecretKey string
)

type Tenants interface {
	SecretKey(ak AccessKey) (SecretKey, error)
	ListAll() ([]AccessKey, error)
}
