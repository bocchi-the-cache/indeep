package fsmdb

import (
	"encoding/binary"
	"errors"

	"github.com/cockroachdb/pebble"
)

type DB struct {
	*pebble.DB
	bin binary.ByteOrder
}

func Open(dirname string, options *pebble.Options) (*DB, error) {
	db, err := pebble.Open(dirname, options)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, bin: binary.LittleEndian}, nil
}

func keyUpperBound(prefix []byte) []byte {
	ub := make([]byte, len(prefix))
	copy(ub, prefix)
	for i := len(ub) - 1; i >= 0; i-- {
		ub[i] = ub[i] + 1
		if ub[i] != 0 {
			return ub[:i+1]
		}
	}
	return nil
}

func prefixIterOptions(prefix string) *pebble.IterOptions {
	bytes := []byte(prefix)
	return &pebble.IterOptions{LowerBound: bytes, UpperBound: keyUpperBound(bytes)}
}

func (d *DB) NewPrefixIter(prefix string) (*pebble.Iterator, error) {
	return d.NewIter(prefixIterOptions(prefix))
}

func (d *DB) Inc(key string) (old uint64, err error) {
	bytesKey := []byte(key)
	newBytes := make([]byte, 8)

	data, closer, err := d.DB.Get(bytesKey)
	if errors.Is(err, pebble.ErrNotFound) {
		d.bin.PutUint64(newBytes, 1)
		err = d.DB.Set(bytesKey, newBytes, nil)
		return
	}
	if err != nil {
		return
	}
	defer func() { _ = closer.Close() }()

	old = d.bin.Uint64(data)
	d.bin.PutUint64(newBytes, old+1)
	err = d.DB.Set(bytesKey, newBytes, nil)
	return
}
