package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	HostScheme = "http"
	RaftScheme = "tcp"
	RootPath   = "/"
)

type Instance interface {
	fmt.Stringer

	URL() *url.URL
	RPC(id RpcID) *url.URL

	json.Marshaler
	json.Unmarshaler
}

type URLInstance struct{ u *url.URL }

func NewURLInstance(scheme, host string) Instance {
	return &URLInstance{u: &url.URL{Scheme: scheme, Host: host}}
}

func (i *URLInstance) String() string { return i.u.String() }
func (i *URLInstance) URL() *url.URL  { return i.u }
func (i *URLInstance) RPC(id RpcID) *url.URL {
	return &url.URL{Scheme: i.u.Scheme, Host: i.u.Host, Path: RootPath + string(id)}
}

func (i *URLInstance) MarshalJSON() ([]byte, error) { return json.Marshal(i.String()) }

func (i *URLInstance) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	i.u = u
	return nil
}
