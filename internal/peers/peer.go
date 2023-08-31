package peers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bocchi-the-cache/indeep/api"
)

const (
	OperationGetMembers        = "/get-members"
	OperationAskLeader         = "/ask-leader"
	OperationLookupMetaService = "/lookup-meta-service"
	OperationAddMetaService    = "/add-meta-service"
	OperationLookupDataService = "/lookup-data-service"
	OperationAddDataService    = "/add-data-service"

	DefaultScheme = "http"
	IDsPrefix     = "/"
	mapSep        = ","
)

var (
	ErrEmptyHosts    = errors.New("empty hosts")
	ErrEmptyIDs      = errors.New("empty peer IDs")
	ErrMapMismatched = errors.New("map mismatched")
)

type peers struct{ m map[api.PeerID]api.Peer }

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
		Scheme: DefaultScheme,
		Host:   strings.Join(hosts, mapSep),
		Path:   IDsPrefix + strings.Join(ids, mapSep),
	}
	return u.String()
}

func DefaultPeers() api.Peers { return new(peers) }

func ParsePeers(rawURL string) (api.Peers, error) {
	m, err := parsePeers(rawURL)
	if err != nil {
		return nil, err
	}
	return &peers{m}, nil
}

func parsePeers(rawURL string) (map[api.PeerID]api.Peer, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, ErrEmptyHosts
	}
	hosts := strings.Split(u.Host, mapSep)

	rawPath := strings.TrimLeft(u.Path, IDsPrefix)
	if rawPath == "" {
		return nil, ErrEmptyIDs
	}
	ids := strings.Split(rawPath, mapSep)

	if len(hosts) != len(ids) {
		return nil, fmt.Errorf("%w: hosts=%v, ids=%v", ErrMapMismatched, hosts, ids)
	}

	m := make(map[api.PeerID]api.Peer)
	for i, host := range hosts {
		id := api.PeerID(ids[i])
		m[id] = &peer{u: &url.URL{Scheme: DefaultScheme, Host: host}}
	}

	return m, nil
}

func (p *peers) IDs() (ret []api.PeerID) {
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

func (p *peers) Lookup(id api.PeerID) api.Peer { return p.m[id] }

func (p *peers) Join(id api.PeerID, peer api.Peer) api.Peers {
	if _, ok := p.m[id]; !ok {
		p.m[id] = peer
	}
	return p
}

func (p *peers) Quit(id api.PeerID) { delete(p.m, id) }

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

	rawPath := strings.TrimLeft(u.Path, IDsPrefix)
	if rawPath == "" {
		return ErrEmptyIDs
	}
	ids := strings.Split(rawPath, mapSep)

	if len(hosts) != len(ids) {
		return fmt.Errorf("%w: hosts=%v, ids=%v", ErrMapMismatched, hosts, ids)
	}

	m := make(map[api.PeerID]api.Peer)
	for i, host := range hosts {
		id := api.PeerID(ids[i])
		m[id] = &peer{u: &url.URL{Scheme: DefaultScheme, Host: host}}
	}
	p.m = m

	return nil
}

type peer struct{ u *url.URL }

func DefaultPeer() api.Peer { return new(peer) }

func ParsePeer(rawURL string) (api.Peer, error) {
	p, err := parsePeer(rawURL)
	if err != nil {
		return nil, err
	}
	return &peer{p}, nil
}

var parsePeer = url.Parse

func (p *peer) String() string { return p.u.String() }

func (p *peer) URL() *url.URL { return &url.URL{Scheme: p.u.Scheme, Host: p.u.Host} }

func (p *peer) Operation(op string) *url.URL {
	u := p.URL()
	u.Path = op
	return u
}

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
