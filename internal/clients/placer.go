package clients

import (
	"net/http"
	"time"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/hyped"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

type PlacerConfig struct {
	HostMap       *api.AddressMap
	ClientTimeout time.Duration
}

type placerClient struct {
	config  *PlacerConfig
	rpc     hyped.RPC
	members api.Peers
	leader  api.Peer
}

func NewPlacer(c *PlacerConfig) (api.Placer, error) {
	cl := &placerClient{
		config:  c,
		rpc:     hyped.NewRPC(&http.Client{Timeout: c.ClientTimeout}),
		members: peers.NewPeers(c.HostMap),
	}
	leader, err := api.AskLeader(cl)
	if err != nil {
		return nil, err
	}
	cl.leader = leader
	return cl, nil
}

func (c *placerClient) GetMembers() api.Peers { return c.members }

func (c *placerClient) AskLeader(p api.Peer) (api.Peer, error) {
	leader := peers.DefaultPeer()
	if err := c.rpc.Get(p, api.RpcMemberAskLeader, leader); err != nil {
		return nil, err
	}
	return leader, nil
}

func (c *placerClient) ListGroups() (ret []api.GroupID, err error) {
	err = c.rpc.Get(c.leader, api.RpcPlacerListGroups, &ret)
	return
}
