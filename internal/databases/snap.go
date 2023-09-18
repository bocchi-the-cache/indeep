package databases

import (
	"io"

	"github.com/hashicorp/raft"
)

func (d *DB) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DB) List() ([]*raft.SnapshotMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DB) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}
