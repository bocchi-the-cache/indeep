package snapmetadb

import (
	"encoding/binary"

	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type DB struct {
	db  *raftboltdb.BoltStore
	bin binary.ByteOrder
}

func Open(path string) (*DB, error) {
	db, err := raftboltdb.New(raftboltdb.Options{Path: path})
	if err != nil {
		return nil, err
	}
	return &DB{db: db, bin: binary.LittleEndian}, nil
}
