package fsmdb

import (
	"encoding/binary"
	"errors"
	"path/filepath"

	"github.com/bocchi-the-cache/indeep/internal/logs"

	"github.com/dgraph-io/badger/v4"
)

const StateDir = "states"

type DB struct {
	*badger.DB
	dataDir string
	bin     binary.ByteOrder
}

func Open(dataDir string) (*DB, error) {
	opts := badger.DefaultOptions(filepath.Join(dataDir, StateDir))
	opts.Logger = logs.Badger()
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, dataDir: dataDir, bin: binary.LittleEndian}, nil
}

func (d *DB) Iter(prefix []byte, withValues bool) *badger.Iterator {
	tx := d.NewTransaction(false)
	opt := badger.IteratorOptions{Prefix: prefix}
	if withValues {
		opt.PrefetchValues = true
		opt.PrefetchSize = badger.DefaultIteratorOptions.PrefetchSize
	}
	return tx.NewIterator(opt)
}

func (d *DB) Uint64(key, dst []byte) (uint64, error) {
	item, err := d.NewTransaction(false).Get(key)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	data, err := item.ValueCopy(dst)
	if err != nil {
		return 0, err
	}

	return d.bin.Uint64(data), nil
}

func (d *DB) Inc(key []byte) (old uint64, err error) {
	dst := make([]byte, 8)

	old, err = d.Uint64(key, dst)
	if err != nil {
		return
	}
	d.bin.PutUint64(dst, old+1)

	tx := d.NewTransaction(true)
	if err = tx.Set(key, dst); err != nil {
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	return
}
