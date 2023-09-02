package peers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
)

const (
	HostScheme = "http"
	RaftScheme = "tcp"

	RootPath = "/"
	mapSep   = ","
)

var (
	ErrEmptyHosts    = errors.New("empty hosts")
	ErrEmptyIDs      = errors.New("empty peer IDs")
	ErrMapMismatched = errors.New("map mismatched")
)

type peers struct {
	scheme string
	m      map[raft.ServerID]api.Peer
}

func RaftPeers() api.Peers { return NewPeers(RaftScheme) }
func HostPeers() api.Peers { return NewPeers(HostScheme) }

func NewPeers(scheme string) api.Peers {
	return &peers{scheme: scheme, m: make(map[raft.ServerID]api.Peer)}
}

func ParsePeers(rawURL string) (api.Peers, error) {
	s, m, err := parsePeers(rawURL)
	if err != nil {
		return nil, err
	}
	return &peers{scheme: s, m: m}, nil
}

func parsePeers(rawURL string) (string, map[raft.ServerID]api.Peer, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", nil, err
	}

	if u.Host == "" {
		return "", nil, ErrEmptyHosts
	}
	hosts := strings.Split(u.Host, mapSep)

	rawPath := strings.TrimLeft(u.Path, RootPath)
	if rawPath == "" {
		return "", nil, ErrEmptyIDs
	}
	ids := strings.Split(rawPath, mapSep)

	if len(hosts) != len(ids) {
		return "", nil, fmt.Errorf("%w: hosts=%v, ids=%v", ErrMapMismatched, hosts, ids)
	}

	m := make(map[raft.ServerID]api.Peer)
	for i, host := range hosts {
		id := raft.ServerID(ids[i])
		m[id] = &peer{u: &url.URL{Scheme: u.Scheme, Host: host}}
	}

	return u.Scheme, m, nil
}

func (p *peers) String() string {
	var (
		ids   []string
		hosts []string
	)
	for id, e := range p.m {
		ids = append(ids, string(id))
		hosts = append(hosts, e.URL().Host)
	}
	u := &url.URL{
		Scheme: p.scheme,
		Host:   strings.Join(hosts, mapSep),
		Path:   RootPath + strings.Join(ids, mapSep),
	}
	return u.String()
}

func (p *peers) IDs() (ret []raft.ServerID) {
	for id := range p.m {
		ret = append(ret, id)
	}
	return
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
			Address:  raft.ServerAddress(peer.URL().Host),
		})
	}
	return
}

func (p *peers) Lookup(id raft.ServerID) api.Peer { return p.m[id] }

func (p *peers) Join(id raft.ServerID, peer api.Peer) api.Peers {
	if _, ok := p.m[id]; !ok {
		p.m[id] = peer
	}
	return p
}

func (p *peers) Quit(id raft.ServerID) { delete(p.m, id) }

func (p *peers) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

func (p *peers) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if u.Host == "" {
		return ErrEmptyHosts
	}
	hosts := strings.Split(u.Host, mapSep)

	rawPath := strings.TrimLeft(u.Path, RootPath)
	if rawPath == "" {
		return ErrEmptyIDs
	}
	ids := strings.Split(rawPath, mapSep)

	if len(hosts) != len(ids) {
		return fmt.Errorf("%w: hosts=%v, ids=%v", ErrMapMismatched, hosts, ids)
	}

	m := make(map[raft.ServerID]api.Peer)
	for i, host := range hosts {
		id := raft.ServerID(ids[i])
		m[id] = &peer{u: &url.URL{Scheme: p.scheme, Host: host}}
	}
	p.m = m

	return nil
}

type peer struct {
	u *url.URL
	s raft.ServerSuffrage
}

func DefaultPeer() api.Peer { return new(peer) }

func Voter(addr raft.ServerAddress) api.Peer {
	return &peer{u: &url.URL{Scheme: RaftScheme, Host: string(addr)}}
}

var parsePeer = url.Parse

func (p *peer) String() string { return p.u.String() }
func (p *peer) URL() *url.URL  { return p.u }
func (p *peer) RPC(id api.RpcID) *url.URL {
	return &url.URL{Scheme: p.u.Scheme, Host: p.u.Host, Path: RootPath + string(id)}
}

func (p *peer) Suffrage() raft.ServerSuffrage { return p.s }

func (p *peer) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

func (p *peer) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	u, err := parsePeer(rawURL)
	if err != nil {
		return err
	}
	p.u = u
	return nil
}

type (
	PeerMux interface {
		HandleFunc(rpc api.RpcID, f http.HandlerFunc) PeerMux
		Build() *http.ServeMux
	}

	peerMux struct {
		p api.Peer
		m *http.ServeMux
	}
)

func Mux(p api.Peer) PeerMux { return &peerMux{p: p, m: http.NewServeMux()} }

func (s *peerMux) HandleFunc(rpc api.RpcID, f http.HandlerFunc) PeerMux {
	s.m.HandleFunc(s.p.RPC(rpc).Path, f)
	return s
}

func (s *peerMux) Build() *http.ServeMux { return s.m }
