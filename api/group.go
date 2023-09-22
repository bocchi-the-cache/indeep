package api

import (
	"io"

	"github.com/hashicorp/raft"
)

type (
	NodeHost string
	GroupID  string
)

type StreamLayerMux interface {
	io.Closer
	NetworkLayer(groupID GroupID) raft.StreamLayer
}
