package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/raft"
)

const (
	HostScheme = "http"
	RaftScheme = "tcp"
	RootPath   = "/"
	HostSep    = ","
)

var (
	ErrAddressEmptyHosts     = errors.New("empty hosts in address")
	ErrAddressEmptyServerIDs = errors.New("empty server IDs in address")
	ErrAddressMapMismatched  = errors.New("map mismatched in address")
)

type RpcID string

type Addresser interface {
	Address() *Address
}

type Address struct{ *url.URL }

var (
	_ = (fmt.Stringer)((*Address)(nil))
	_ = (json.Marshaler)((*Address)(nil))
	_ = (json.Unmarshaler)((*Address)(nil))
)

func NewAddress(scheme, host string) *Address {
	return &Address{&url.URL{Scheme: scheme, Host: host}}
}

func (a *Address) RPC(id RpcID) *Address {
	return &Address{&url.URL{Scheme: a.Scheme, Host: a.Host, Path: RootPath + string(id)}}
}

func (a *Address) MarshalJSON() ([]byte, error) { return json.Marshal(a.String()) }

func (a *Address) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	a.URL = u
	return nil
}

type AddressList struct {
	scheme string
	hosts  []string
}

var (
	_ = (fmt.Stringer)((*AddressList)(nil))
	_ = (json.Marshaler)((*AddressList)(nil))
	_ = (json.Unmarshaler)((*AddressList)(nil))
)

func (a *AddressList) String() string {
	return (&url.URL{Scheme: a.scheme, Host: strings.Join(a.hosts, HostSep)}).String()
}

func (a *AddressList) MarshalJSON() ([]byte, error) { return json.Marshal(a.String()) }

func (a *AddressList) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	scheme, hosts, err := parseAddressList(rawURL)
	if err != nil {
		return err
	}
	a.scheme = scheme
	a.hosts = hosts
	return nil
}

func parseAddressList(rawURL string) (scheme string, hosts []string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	if u.Host == "" {
		err = ErrAddressEmptyHosts
		return
	}
	scheme = u.Scheme
	hosts = strings.Split(u.Host, HostSep)
	return
}

type AddressMap struct {
	scheme string
	hosts  map[raft.ServerID]string
}

func parseAddressMap(rawURL string) (scheme string, hosts map[raft.ServerID]string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	if u.Host == "" {
		err = ErrAddressEmptyHosts
		return
	}
	hostList := strings.Split(u.Host, HostSep)

	rawPath := strings.TrimLeft(u.Path, RootPath)
	if rawPath == "" {
		err = ErrAddressEmptyServerIDs
		return
	}
	ids := strings.Split(rawPath, HostSep)

	if len(hostList) != len(ids) {
		err = fmt.Errorf("%w: hostList=%v, ids=%v", ErrAddressMapMismatched, hostList, ids)
		return
	}

	scheme = u.Scheme
	hosts = make(map[raft.ServerID]string)
	for i, host := range hostList {
		id := raft.ServerID(ids[i])
		hosts[id] = host
	}

	return
}

func (p *AddressMap) String() string {
	var (
		ids   []string
		hosts []string
	)
	for id, host := range p.hosts {
		ids = append(ids, string(id))
		hosts = append(hosts, host)
	}
	u := &url.URL{
		Scheme: p.scheme,
		Host:   strings.Join(hosts, HostSep),
		Path:   RootPath + strings.Join(ids, HostSep),
	}
	return u.String()
}

func (p *AddressMap) Lookup(id raft.ServerID) string { return p.hosts[id] }

func (p *AddressMap) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

func (p *AddressMap) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}

	scheme, hosts, err := parseAddressMap(rawURL)
	if err != nil {
		return err
	}

	p.scheme = scheme
	p.hosts = hosts
	return nil
}
