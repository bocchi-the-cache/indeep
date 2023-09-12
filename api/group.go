package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
)

type (
	NodeHost string
	GroupID  string
)

type Multipeer interface {
	Addresser
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
