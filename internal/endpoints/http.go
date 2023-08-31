package endpoints

import (
	"encoding/json"
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

	httpListSep = ","
	httpScheme  = "http"
)

type endpointList struct{ es []api.Endpoint }

func DefaultEndpointList() api.EndpointList { return new(endpointList) }

func (h *endpointList) Endpoints() []api.Endpoint { return h.es }

func (h *endpointList) MarshalJSON() ([]byte, error) {
	var hosts []string
	for _, e := range h.es {
		hosts = append(hosts, e.URL().Host)
	}
	u := &url.URL{Scheme: httpScheme, Host: strings.Join(hosts, httpListSep)}
	return json.Marshal(u.String())
}

func (h *endpointList) UnmarshalJSON(bytes []byte) error {
	var rawURL string
	if err := json.Unmarshal(bytes, &rawURL); err != nil {
		return err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	for _, host := range strings.Split(u.Host, ",") {
		h.es = append(h.es, NewEndpoint(&url.URL{Scheme: httpScheme, Host: host}))
	}
	return nil
}

type endpoint struct{ u *url.URL }

func DefaultEndpoint() api.Endpoint       { return new(endpoint) }
func NewEndpoint(u *url.URL) api.Endpoint { return &endpoint{u} }

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
