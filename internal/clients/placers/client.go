package placers

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
)

type client struct {
	h *http.Client
}

func (c *client) LookupMetaClient(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddMetaServer() error {
	//TODO implement me
	panic("implement me")
}

func (c *client) LookupDataClient(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddDataServer() error {
	//TODO implement me
	panic("implement me")
}
