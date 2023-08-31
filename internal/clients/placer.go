package clients

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/endpoints"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
)

type PlacerConfig struct {
	EndpointMap api.EndpointMap

	httpClient *http.Client
}

type placerClient struct {
	h *http.Client

	members api.EndpointMap
	leader  api.Endpoint
}

func NewPlacer(c *PlacerConfig) (api.Placer, error) {
	cl := &placerClient{
		h:       c.httpClient,
		members: c.EndpointMap,
	}
	leader, err := api.AskLeader(cl)
	if err != nil {
		return nil, err
	}
	cl.leader = leader
	return cl, nil
}

func (c *placerClient) Members() (ret []api.Endpoint) {
	for _, e := range c.members.Endpoints() {
		ret = append(ret, e)
	}
	return
}

func (c *placerClient) Leader(e api.Endpoint) (api.Endpoint, error) {
	resp, err := c.h.Get(e.Operation(endpoints.OperationAskLeader).String())
	if err != nil {
		return nil, err
	}

	leader := endpoints.DefaultEndpoint()
	if err := jsonhttp.Unmarshal(resp.Body, leader); err != nil {
		return nil, err
	}

	return leader, nil
}

func (c *placerClient) LookupMetaService(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) AddMetaService() error {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) LookupDataService(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *placerClient) AddDataService() error {
	//TODO implement me
	panic("implement me")
}
