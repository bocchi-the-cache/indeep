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

func (id RpcID) Path() string { return RootPath + string(id) }

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
	return &Address{&url.URL{Scheme: a.Scheme, Host: a.Host, Path: id.Path()}}
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
	hosts  []raft.ServerAddress
}

var (
	_ = (fmt.Stringer)((*AddressList)(nil))
	_ = (json.Marshaler)((*AddressList)(nil))
	_ = (json.Unmarshaler)((*AddressList)(nil))
)

func (l *AddressList) String() string {
	var rawHosts []string
	for _, host := range l.hosts {
		rawHosts = append(rawHosts, string(host))
	}
	return (&url.URL{Scheme: l.scheme, Host: strings.Join(rawHosts, HostSep)}).String()
}

func (l *AddressList) MarshalJSON() ([]byte, error) { return json.Marshal(l.String()) }

func (l *AddressList) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	scheme, hosts, err := parseAddressList(rawURL)
	if err != nil {
		return err
	}
	l.scheme = scheme
	l.hosts = hosts
	return nil
}

func parseAddressList(rawURL string) (scheme string, hosts []raft.ServerAddress, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	if u.Host == "" {
		err = ErrAddressEmptyHosts
		return
	}
	scheme = u.Scheme
	for _, s := range strings.Split(u.Host, HostSep) {
		hosts = append(hosts, raft.ServerAddress(s))
	}
	return
}

type AddressMap struct {
	scheme string
	hosts  map[raft.ServerID]raft.ServerAddress
}

func NewAddressMap(scheme string) *AddressMap {
	return &AddressMap{scheme: scheme, hosts: make(map[raft.ServerID]raft.ServerAddress)}
}

func (m *AddressMap) Scheme() string { return m.scheme }

func (m *AddressMap) Addresses() map[raft.ServerID]raft.ServerAddress {
	return m.hosts
}

func (m *AddressMap) Join(id raft.ServerID, host string) *AddressMap {
	if _, ok := m.hosts[id]; !ok {
		m.hosts[id] = raft.ServerAddress(host)
	}
	return m
}

func ParseAddressMap(rawURL string) (*AddressMap, error) {
	scheme, hosts, err := parseAddressMap(rawURL)
	if err != nil {
		return nil, err
	}
	return &AddressMap{scheme, hosts}, nil
}

func parseAddressMap(rawURL string) (scheme string, hosts map[raft.ServerID]raft.ServerAddress, err error) {
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
	hosts = make(map[raft.ServerID]raft.ServerAddress)
	for i, host := range hostList {
		hosts[raft.ServerID(ids[i])] = raft.ServerAddress(host)
	}

	return
}

func (m *AddressMap) String() string {
	var (
		ids   []string
		hosts []string
	)
	for id, host := range m.hosts {
		ids = append(ids, string(id))
		hosts = append(hosts, string(host))
	}
	u := &url.URL{
		Scheme: m.scheme,
		Host:   strings.Join(hosts, HostSep),
		Path:   RootPath + strings.Join(ids, HostSep),
	}
	return u.String()
}

func (m *AddressMap) MarshalJSON() ([]byte, error) { return json.Marshal(m.String()) }

func (m *AddressMap) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}

	scheme, hosts, err := parseAddressMap(rawURL)
	if err != nil {
		return err
	}

	m.scheme = scheme
	m.hosts = hosts
	return nil
}
