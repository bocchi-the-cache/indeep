package peers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
)

type peers struct {
	m map[raft.ServerID]api.Peer
}

func NewPeers(a *api.AddressMap) (api.Peers, error) {
	m := make(map[raft.ServerID]api.Peer)
	for id, host := range a.Addresses() {
		m[id] = &peer{id: id, addr: api.NewAddress(a.Scheme(), string(host))}
	}
	if len(m) == 0 {
		return nil, api.ErrEmptyPeers
	}
	return &peers{m: m}, nil
}

func (p *peers) Peers() (ret []api.Peer) {
	for _, peer := range p.m {
		ret = append(ret, peer)
	}
	return
}

func (p *peers) Configuration() (c raft.Configuration) {
	for id, peer := range p.m {
		c.Servers = append(c.Servers, raft.Server{
			Suffrage: peer.Suffrage(),
			ID:       id,
			Address:  raft.ServerAddress(peer.Address().Host),
		})
	}
	return
}

func (p *peers) Lookup(id raft.ServerID) (api.Peer, error) {
	peer, ok := p.m[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", api.ErrPeerUnknown, id)
	}
	return peer, nil
}

type peerInfo struct {
	ID      raft.ServerID
	Address *api.Address
}

type peer struct {
	id   raft.ServerID
	addr *api.Address
	s    raft.ServerSuffrage
}

func (p *peer) Address() *api.Address         { return p.addr }
func (p *peer) ID() raft.ServerID             { return p.id }
func (p *peer) Suffrage() raft.ServerSuffrage { return p.s }

func (p *peer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&peerInfo{ID: p.id, Address: p.addr})
}

func (p *peer) UnmarshalJSON(bytes []byte) error {
	var info peerInfo
	if err := json.Unmarshal(bytes, &info); err != nil {
		return err
	}
	p.id = info.ID
	p.addr = info.Address
	return nil
}

type (
	PeerServeMux interface {
		HandleFunc(rpc api.RpcID, f http.HandlerFunc) PeerServeMux
		Build() *http.ServeMux
	}

	peerServeMux struct {
		p api.Peer
		m *http.ServeMux
	}
)

func ServeMux(p api.Peer) PeerServeMux { return &peerServeMux{p: p, m: http.NewServeMux()} }

func (s *peerServeMux) HandleFunc(rpc api.RpcID, f http.HandlerFunc) PeerServeMux {
	s.m.HandleFunc(s.p.Address().RPC(rpc).Path, f)
	return s
}

func (s *peerServeMux) Build() *http.ServeMux { return s.m }
