package clients

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/jsonutl"
)

const (
	httpPathIsLeader = "/get"
)

type placerClient struct {
	h *http.Client

	members []api.Endpoint
	leader  api.Endpoint
}

func (c *placerClient) Members() []api.Endpoint { return c.members }

func (c *placerClient) IsLeader(e api.Endpoint) (bool, error) {
	// FIXME: Ask one member could know the leader of this cluster, it's unnecessary to ask all.
	u := e.URL()
	u.Path = httpPathIsLeader
	resp, err := c.h.Get(u.String())
	if err != nil {
		return false, err
	}

	var isLeader bool
	if err := jsonutl.UnmarshalBody(resp.Body, &isLeader); err != nil {
		return false, err
	}

	return isLeader, nil
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
