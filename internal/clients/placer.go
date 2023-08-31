package clients

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

type PlacerConfig struct {
	Peers    api.Peers
	RawPeers string

	httpClient *http.Client
}

type placerClient struct {
	h *http.Client

	members api.Peers
	leader  api.Peer
}

func NewPlacer(c *PlacerConfig) (api.Placer, error) {
	cl := &placerClient{
		h:       c.httpClient,
		members: c.Peers,
	}
	leader, err := api.AskLeader(cl)
	if err != nil {
		return nil, err
	}
	cl.leader = leader
	return cl, nil
}

func (c *placerClient) Members() api.Peers { return c.members }

func (c *placerClient) Leader(e api.Peer) (api.Peer, error) {
	resp, err := c.h.Get(e.Operation(peers.OperationAskLeader).String())
	if err != nil {
		return nil, err
	}

	leader := peers.DefaultPeer()
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
