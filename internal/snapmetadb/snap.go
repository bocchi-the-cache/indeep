package snapmetadb

import (
	"io"

	"github.com/hashicorp/raft"
)

const (
	SnapshotVersionKey            = "SnapshotVersion"
	SnapshotIndexKey              = "SnapshotIndex"
	SnapshotTermKey               = "SnapshotTerm"
	SnapshotConfigurationKey      = "SnapshotConfiguration"
	SnapshotConfigurationIndexKey = "SnapshotConfigurationIndex"
)

func (d *DB) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
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
	if err != nil {
		return nil, err
	}
	return []*raft.SnapshotMeta{m}, nil
}

func (d *DB) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	m, err := d.Get()
	if err != nil {
		return nil, nil, err
	}
	return m, nil, nil
}

func (*DB) Write([]byte) (int, error) { return 0, nil }
func (*DB) Close() error              { return nil }
func (d *DB) ID() string              { return string(d.id) }
func (*DB) Cancel() error             { return nil }
