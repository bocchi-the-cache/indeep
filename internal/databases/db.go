package databases

import (
	"github.com/cockroachdb/pebble"
)

type DB struct {
	db *pebble.DB
}

func Open(path string) (*DB, error) {
	db, err := pebble.Open(path, new(pebble.Options))
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
