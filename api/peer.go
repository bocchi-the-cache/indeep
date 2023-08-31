package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type PeerID string

type Peers interface {
	fmt.Stringer

	IDs() []PeerID
	Peers() []Peer

	Lookup(id PeerID) Peer
	Join(id PeerID, peer Peer) Peers
	Quit(id PeerID)

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
