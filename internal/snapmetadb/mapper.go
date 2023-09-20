package snapmetadb

import (
	"encoding/json"

	"github.com/hashicorp/raft"
)

type Mapper interface {
	Save(meta raft.SnapshotMeta) error
	Get() (*raft.SnapshotMeta, error)
}

func (d *DB) Save(meta raft.SnapshotMeta) error {
	{
		bytes := make([]byte, 8)

		d.bin.PutUint64(bytes, uint64(meta.Version))
		if err := d.db.Set([]byte(SnapshotVersionKey), bytes); err != nil {
			return err
		}

		d.bin.PutUint64(bytes, meta.Index)
		if err := d.db.Set([]byte(SnapshotIndexKey), bytes); err != nil {
			return err
		}

		d.bin.PutUint64(bytes, meta.Term)
		if err := d.db.Set([]byte(SnapshotTermKey), bytes); err != nil {
			return err
		}

		d.bin.PutUint64(bytes, meta.ConfigurationIndex)
		if err := d.db.Set([]byte(SnapshotConfigurationIndexKey), bytes); err != nil {
			return err
		}
	}

	{
		bytes, err := json.Marshal(meta.Configuration)
		if err != nil {
			return err
		}
		if err := d.db.Set([]byte(SnapshotConfigurationKey), bytes); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) Get() (*raft.SnapshotMeta, error) {
	ret := new(raft.SnapshotMeta)

	{
		var (
			bytes []byte
			err   error
		)

		bytes, err = d.db.Get([]byte(SnapshotVersionKey))
		if err != nil {
			return nil, err
		}
		ret.Version = raft.SnapshotVersion(d.bin.Uint64(bytes))

		bytes, err = d.db.Get([]byte(SnapshotIndexKey))
		if err != nil {
			return nil, err
		}
		ret.Index = d.bin.Uint64(bytes)

		bytes, err = d.db.Get([]byte(SnapshotTermKey))
		if err != nil {
			return nil, err
		}
		ret.Term = d.bin.Uint64(bytes)

		bytes, err = d.db.Get([]byte(SnapshotConfigurationIndexKey))
		if err != nil {
			return nil, err
		}
		ret.ConfigurationIndex = d.bin.Uint64(bytes)
	}

	{
		bytes, err := d.db.Get([]byte(SnapshotConfigurationKey))
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bytes, &ret.Configuration); err != nil {
			return nil, err
		}
	}

	return ret, nil
}
