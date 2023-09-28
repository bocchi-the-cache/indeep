package fsmdb

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"

	"github.com/bocchi-the-cache/indeep/internal/logs"
)

const StateDir = "states"

type DB struct {
	*badger.DB
	dataDir string
	bin     binary.ByteOrder
}

func Open(dataDir string) (*DB, error) {
	fsmDB := &DB{dataDir: dataDir, bin: binary.LittleEndian}

	stateDir := filepath.Join(dataDir, StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return nil, err
	}
	snapDir := fsmDB.snapDir()
	if err := os.MkdirAll(snapDir, 0755); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(stateDir)
	opts.Logger = logs.Badger()

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	fsmDB.DB = db

	return fsmDB, nil
}

func (d *DB) View(f func(tx *Tx) error) error {
	tx := &Tx{Txn: d.NewTransaction(false), bin: d.bin}
	defer tx.Discard()
	return f(tx)
}

func (d *DB) Update(f func(tx *Tx) error) error {
	tx := &Tx{Txn: d.NewTransaction(true), bin: d.bin}
	defer tx.Discard()
	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit()
}

type Tx struct {
	*badger.Txn
	bin binary.ByteOrder
}

func (tx *Tx) Iter(prefix []byte, withValues bool) *badger.Iterator {
	opt := badger.IteratorOptions{Prefix: prefix}
	if withValues {
		opt.PrefetchValues = true
		opt.PrefetchSize = badger.DefaultIteratorOptions.PrefetchSize
	}
	return tx.NewIterator(opt)
}

func (tx *Tx) Uint64(key, dst []byte) (uint64, error) {
	item, err := tx.Get(key)
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

	return tx.bin.Uint64(data), nil
}

func (tx *Tx) Inc(key []byte) (old uint64, err error) {
	dst := make([]byte, 8)

	old, err = tx.Uint64(key, dst)
	if err != nil {
		return
	}
	tx.bin.PutUint64(dst, old+1)

	err = tx.Set(key, dst)
	return
}
