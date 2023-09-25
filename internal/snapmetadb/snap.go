package snapmetadb

import (
	"errors"
	"io"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// FIXME: Dump the FSM DB snapshot, not just record the snapshot meta.

const (
	SnapshotVersionKey            = "SnapshotVersion"
	SnapshotIndexKey              = "SnapshotIndex"
	SnapshotTermKey               = "SnapshotTerm"
	SnapshotConfigurationKey      = "SnapshotConfiguration"
	SnapshotConfigurationIndexKey = "SnapshotConfigurationIndex"
)

func (d *DB) Create(
	version raft.SnapshotVersion,
	index, term uint64,
	configuration raft.Configuration,
	configurationIndex uint64,
	_ raft.Transport,
) (raft.SnapshotSink, error) {
	if err := d.Save(raft.SnapshotMeta{
		Version:            version,
		Index:              index,
		Term:               term,
		Configuration:      configuration,
		ConfigurationIndex: configurationIndex,
	}); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DB) List() ([]*raft.SnapshotMeta, error) {
	m, err := d.Get()
	if errors.Is(err, raftboltdb.ErrKeyNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []*raft.SnapshotMeta{m}, nil
}

func (d *DB) Open(string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	m, err := d.Get()
	if err != nil {
		return nil, nil, err
	}
	return m, nil, nil
}

func (*DB) Write([]byte) (_ int, _ error) { return }
func (*DB) Close() (_ error)              { return }
func (*DB) ID() (_ string)                { return }
func (*DB) Cancel() (_ error)             { return }
