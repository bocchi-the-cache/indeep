package api

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/raft"
)

var (
	ErrEmptyPeers  = errors.New("empty peers")
	ErrPeerUnknown = errors.New("peer unknown")
)

type Peers interface {
	Peers() []Peer
	Configuration() raft.Configuration
	Lookup(id raft.ServerID) (Peer, error)
}

type Peer interface {
	Addresser

	ID() raft.ServerID
	Suffrage() raft.ServerSuffrage

	json.Marshaler
	json.Unmarshaler
}
