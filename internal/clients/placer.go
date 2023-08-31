package clients

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/endpoints"
	"github.com/bocchi-the-cache/indeep/internal/jsonutl"
)

const (
	operationAskLeader = "/leader"
)

type placerClient struct {
	h *http.Client

	members []api.Endpoint
	leader  api.Endpoint
}

func (c *placerClient) Members() []api.Endpoint { return c.members }

func (c *placerClient) Leader(e api.Endpoint) (api.Endpoint, error) {
	resp, err := c.h.Get(e.Operation(operationAskLeader).String())
	if err != nil {
		return nil, err
	}

	leader := endpoints.EmptyHttpEndpoint()
	if err := jsonutl.UnmarshalBody(resp.Body, leader); err != nil {
		return nil, err
	}

	return leader, nil
}

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
