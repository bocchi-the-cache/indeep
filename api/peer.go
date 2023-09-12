package api

import (
	"encoding/json"
	"fmt"

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
	Addresser
	Suffrage() raft.ServerSuffrage
}

type PeerInfo struct {
	ID   raft.ServerID
	Peer Peer
}
