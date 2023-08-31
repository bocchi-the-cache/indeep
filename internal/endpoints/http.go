package endpoints

import (
	"encoding/json"
	"net/url"
)

type httpEndpoint struct{ u *url.URL }

func (e *httpEndpoint) String() string { return e.u.String() }
func (e *httpEndpoint) URL() *url.URL {
	u := *e.u
	return &u
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
