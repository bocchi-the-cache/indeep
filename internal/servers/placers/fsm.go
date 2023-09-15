package placers

import (
	"io"

	"github.com/hashicorp/raft"
)

func (s *placerServer) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}
