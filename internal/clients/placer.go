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
	Peers         api.Peers
	ClientTimeout time.Duration
}

type placerClient struct {
	config   *PlacerConfig
	rpc      hyped.RPC
	members  api.Peers
	leaderID raft.ServerID
	leader   api.Peer
}

func NewPlacer(c *PlacerConfig) (api.Placer, error) {
	cl := &placerClient{
		config:  c,
		rpc:     hyped.NewRPC(&http.Client{Timeout: c.ClientTimeout}),
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

func (c *placerClient) AskLeader(p api.Peer) (*api.PeerInfo, error) {
	info := &api.PeerInfo{Peer: peers.DefaultPeer()}
	if err := c.rpc.Get(p, api.RpcAskLeader, info); err != nil {
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
