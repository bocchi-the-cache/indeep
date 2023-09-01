package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/raft"
)

type Peers interface {
	fmt.Stringer

	IDs() []raft.ServerID
	Peers() []Peer
	Configuration() raft.Configuration

	Lookup(id raft.ServerID) Peer
	Join(id raft.ServerID, peer Peer) Peers
	Quit(id raft.ServerID)

	json.Marshaler
	json.Unmarshaler
}

type Peer interface {
	fmt.Stringer

	URL() *url.URL
	Operation(op string) *url.URL

	Suffrage() raft.ServerSuffrage

	json.Marshaler
	json.Unmarshaler
}
