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

	mapSep        = ","
	pathPrefix    = "/"
	defaultScheme = "http"
)

var (
	ErrEmptyHosts    = errors.New("empty hosts")
	ErrEmptyIDs      = errors.New("empty peer IDs")
	ErrMapMismatched = errors.New("map mismatched")
)

type peers struct{ m map[api.PeerID]api.Peer }

func DefaultPeers() api.Peers { return new(peers) }

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

func (p *peers) MarshalJSON() ([]byte, error) {
	var (
		ids   []string
		hosts []string
	)
	for id, e := range p.m {
		ids = append(ids, string(id))
		hosts = append(hosts, e.URL().Host)
	}
	u := &url.URL{
		Scheme: defaultScheme,
		Host:   strings.Join(hosts, mapSep),
		Path:   pathPrefix + strings.Join(ids, mapSep),
	}
	return json.Marshal(u.String())
}

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

	rawPath := strings.TrimLeft(u.Path, pathPrefix)
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
		m[id] = NewPeer(host)
	}
	p.m = m

	return nil
}

type peer struct{ u *url.URL }

func DefaultPeer() api.Peer        { return new(peer) }
func NewPeer(host string) api.Peer { return &peer{u: &url.URL{Scheme: defaultScheme, Host: host}} }

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
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	p.u = u
	return nil
}
