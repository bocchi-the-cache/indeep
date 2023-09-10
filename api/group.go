package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
)

type (
	NodeAddress string
	GroupID     string
)

type Multipeer interface {
	Instance
}

type Multipeers interface {
	fmt.Stringer

	Join(mp Multipeer)

	json.Marshaler
	json.Unmarshaler
}

type StreamLayerMux interface {
	io.Closer
	NetworkLayer(groupID GroupID) raft.StreamLayer
}
