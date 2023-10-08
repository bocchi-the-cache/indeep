package clients

import (
	"net/http"
	"time"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/peers"
	"github.com/bocchi-the-cache/indeep/internal/utils/hyped"
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
	members, err := peers.NewPeers(c.HostMap)
	if err != nil {
		return nil, err
	}
	cl := &placerClient{
		config:  c,
		rpc:     hyped.NewRPC(&http.Client{Timeout: c.ClientTimeout}),
		members: members,
	}
	leaderID, err := cl.AskLeaderID()
	if err != nil {
		return nil, err
	}
	leader, err := cl.members.Lookup(*leaderID)
	if err != nil {
		return nil, err
	}
	cl.leader = leader
	return cl, nil
}

func (*placerClient) CheckLeader() error { return nil }

func (c *placerClient) AskLeaderID() (*raft.ServerID, error) {
	var leaderID raft.ServerID
	if err := c.rpc.Get(c.members.Peers()[0], api.RpcMemberAskLeaderID, &leaderID); err != nil {
		return nil, err
	}
	return &leaderID, nil
}

func (c *placerClient) ListGroups() (*[]api.GroupID, error) {
	var ret []api.GroupID
	if err := c.rpc.Get(c.leader, api.RpcPlacerListGroups, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *placerClient) GenerateGroup() (*api.GroupID, error) {
	var groupID api.GroupID
	if err := c.rpc.Get(c.leader, api.RpcPlacerGenerateGroup, &groupID); err != nil {
		return nil, err
	}
	return &groupID, nil
}
