package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type PeerID string

type Peers interface {
	IDs() []PeerID
	Peers() []Peer
	Lookup(id PeerID) Peer

	json.Marshaler
	json.Unmarshaler
}

type Peer interface {
	fmt.Stringer

	URL() *url.URL
	Operation(op string) *url.URL

	json.Marshaler
	json.Unmarshaler
}
