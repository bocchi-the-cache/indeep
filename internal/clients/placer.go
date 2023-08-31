package clients

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
)

type placerClient struct {
	h *http.Client

	members []api.Endpoint
	leader  api.Endpoint
}

func (c *placerClient) Members() []api.Endpoint {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) Leader() api.Endpoint { return c.leader }

func (c *placerClient) LookupMetaClient(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) AddMetaServer() error {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) LookupDataClient(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) AddDataServer() error {
	//TODO implement me
	panic("implement me")
}
