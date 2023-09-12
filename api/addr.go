package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	HostScheme = "http"
	RaftScheme = "tcp"
	RootPath   = "/"
	HostSep    = ","
)

var (
	ErrAddressListEmptyHosts = errors.New("empty hosts in address list")
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
		err = ErrAddressListEmptyHosts
		return
	}
	scheme = u.Scheme
	hosts = strings.Split(u.Host, HostSep)
	return
}
