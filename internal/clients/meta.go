package clients

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
)

type metaClient struct {
	h *http.Client
}

func (c *metaClient) Lookup(key api.MetaKey) (api.MetaPartition, error) {
	//TODO implement me
	panic("implement me")
}
