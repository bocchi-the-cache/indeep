package endpoints

import (
	"encoding/json"
	"net/url"

	"github.com/bocchi-the-cache/indeep/api"
)

type httpEndpoint struct{ u *url.URL }

func EmptyHttpEndpoint() api.Endpoint { return new(httpEndpoint) }

func (e *httpEndpoint) String() string { return e.u.String() }

func (e *httpEndpoint) URL() *url.URL { return &url.URL{Scheme: e.u.Scheme, Host: e.u.Host} }
func (e *httpEndpoint) Operation(op string) *url.URL {
	u := e.URL()
	u.Path = op
	return u
}

func (e *httpEndpoint) MarshalJSON() ([]byte, error) { return json.Marshal(e.String()) }
func (e *httpEndpoint) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	e.u = u
	return nil
}
