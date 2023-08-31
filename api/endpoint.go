package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type EndpointList interface {
	Endpoints() []Endpoint

	json.Marshaler
	json.Unmarshaler
}

type Endpoint interface {
	fmt.Stringer

	URL() *url.URL
	Operation(op string) *url.URL

	json.Marshaler
	json.Unmarshaler
}
