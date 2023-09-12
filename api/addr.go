package api

import (
	"encoding/json"
	"net/url"
)

const (
	HostScheme = "http"
	RaftScheme = "tcp"
	RootPath   = "/"
)

type RpcID string

type Addresser interface {
	Address() *Address
}

type Address struct{ *url.URL }

func NewAddress(scheme, host string) *Address {
	return &Address{&url.URL{Scheme: scheme, Host: host}}
}

func (i *Address) RPC(id RpcID) *Address {
	return &Address{&url.URL{Scheme: i.Scheme, Host: i.Host, Path: RootPath + string(id)}}
}

func (i *Address) MarshalJSON() ([]byte, error) { return json.Marshal(i.String()) }

func (i *Address) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	i.URL = u
	return nil
}
