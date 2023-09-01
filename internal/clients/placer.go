package clients

import (
	"net/http"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

type PlacerConfig struct {
	Peers api.Peers

	httpClient *http.Client
}

type placerClient struct {
	h *http.Client

	members  api.Peers
	leaderID raft.ServerID
	leader   api.Peer
}

func NewPlacer(c *PlacerConfig) (api.Placer, error) {
	cl := &placerClient{
		h:       c.httpClient,
		members: c.Peers,
	}
	info, err := api.AskLeader(cl)
	if err != nil {
		return nil, err
	}
	cl.leaderID = info.ID
	cl.leader = info.Peer
	return cl, nil
}

func (c *placerClient) GetMembers() api.Peers { return c.members }

func (c *placerClient) AskLeader(e api.Peer) (*api.PeerInfo, error) {
	resp, err := c.h.Get(e.RPC(api.RpcAskLeader).String())
	if err != nil {
		return nil, err
	}

	info := &api.PeerInfo{Peer: peers.DefaultPeer()}
	if err := jsonhttp.Unmarshal(resp.Body, info); err != nil {
		return nil, err
	}

	return info, nil
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
