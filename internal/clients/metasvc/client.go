package metasvc

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
)

type client struct {
	h *http.Client
}

func (c *client) Lookup(key api.MetaKey) (api.MetaPartition, error) {
	//TODO implement me
	panic("implement me")
}
