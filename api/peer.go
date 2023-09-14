package api

import (
	"encoding/json"

	"github.com/hashicorp/raft"
)

type Peers interface {
	Peers() []Peer
	Configuration() raft.Configuration
	Lookup(id raft.ServerID) Peer
}

type Peer interface {
	Addresser

	ID() raft.ServerID
	Suffrage() raft.ServerSuffrage

	json.Marshaler
	json.Unmarshaler
}
