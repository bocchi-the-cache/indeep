package endpoints

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
	ErrEmptyIDs      = errors.New("empty endpoint IDs")
	ErrMapMismatched = errors.New("map mismatched")
)

type endpointMap struct{ m map[string]api.Endpoint }

func DefaultEndpointMap() api.EndpointMap { return new(endpointMap) }

func (h *endpointMap) Endpoints() map[string]api.Endpoint { return h.m }

func (h *endpointMap) MarshalJSON() ([]byte, error) {
	var (
		ids   []string
		hosts []string
	)
	for id, e := range h.m {
		ids = append(ids, id)
		hosts = append(hosts, e.URL().Host)
	}
	u := &url.URL{
		Scheme: defaultScheme,
		Host:   strings.Join(hosts, mapSep),
		Path:   pathPrefix + strings.Join(ids, mapSep),
	}
	return json.Marshal(u.String())
}

func (h *endpointMap) UnmarshalJSON(bytes []byte) error {
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

	m := make(map[string]api.Endpoint)
	for i, host := range hosts {
		id := ids[i]
		m[id] = NewEndpoint(id, &url.URL{Scheme: defaultScheme, Host: host})
	}
	h.m = m

	return nil
}

type endpoint struct {
	id string
	u  *url.URL
}

func DefaultEndpoint() api.Endpoint                  { return new(endpoint) }
func NewEndpoint(id string, u *url.URL) api.Endpoint { return &endpoint{id, u} }

func (e *endpoint) ID() string     { return e.id }
func (e *endpoint) String() string { return e.u.String() }

func (e *endpoint) URL() *url.URL { return &url.URL{Scheme: e.u.Scheme, Host: e.u.Host} }
func (e *endpoint) Operation(op string) *url.URL {
	u := e.URL()
	u.Path = op
	return u
}

func (e *endpoint) MarshalJSON() ([]byte, error) { return json.Marshal(e.String()) }
func (e *endpoint) UnmarshalJSON(bytes []byte) error {
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
