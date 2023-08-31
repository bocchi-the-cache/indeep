package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type EndpointMap interface {
	Endpoints() map[string]Endpoint

	json.Marshaler
	json.Unmarshaler
}

type Endpoint interface {
	ID() string
	fmt.Stringer

	URL() *url.URL
	Operation(op string) *url.URL

	json.Marshaler
	json.Unmarshaler
}
