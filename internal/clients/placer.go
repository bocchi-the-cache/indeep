package clients

import (
	"net/http"
	"time"

	"github.com/hashicorp/raft"

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
	leaderID, err := api.AskLeaderID(cl)
	if err != nil {
		return nil, err
	}
	leader, err := cl.members.Lookup(leaderID)
	if err != nil {
		return nil, err
	}
	cl.leader = leader
	return cl, nil
}

func (c *placerClient) Peers() api.Peers { return c.members }

func (c *placerClient) AskLeaderID(p api.Peer) (leaderID raft.ServerID, err error) {
	err = c.rpc.Get(p, api.RpcMemberAskLeader, &leaderID)
	return
}

func (c *placerClient) ListGroups() (ret []api.GroupID, err error) {
	err = c.rpc.Get(c.leader, api.RpcPlacerListGroups, &ret)
	return
}
