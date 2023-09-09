package api

import (
	"io"

	"github.com/hashicorp/raft"
)

type (
	NodeAddress string
	GroupID     string
)

type StreamLayerMux interface {
	io.Closer
	NetworkLayer(groupID GroupID) raft.StreamLayer
}
