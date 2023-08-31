package endpoints

import (
	"net/url"

	"github.com/bocchi-the-cache/indeep/api"
)

type httpEndpoint struct{ u *url.URL }

func (e *httpEndpoint) String() string { return e.u.String() }
func (e *httpEndpoint) URL() *url.URL  { return e.u }

func NewHttpEndpoint(rawURL string) (api.Endpoint, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &httpEndpoint{u}, nil
}
