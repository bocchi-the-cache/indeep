package fsmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hashicorp/raft"
)

const (
	SnapDir        = "snaps"
	MetaFileSuffix = ".meta.json"
	SnapFileSuffix = ".badger.bak"
	TimeFileSuffix = ".badger.time"
)

func SnapshotID(index uint64) string { return fmt.Sprintf("snapshot%06d", index) }

func (d *DB) snapDir() string               { return filepath.Join(d.dataDir, SnapDir) }
func (d *DB) metaFilePath(id string) string { return filepath.Join(d.snapDir(), id+MetaFileSuffix) }
func (d *DB) snapFilePath(id string) string { return filepath.Join(d.snapDir(), id+SnapFileSuffix) }
func (d *DB) timeFilePath(id string) string { return filepath.Join(d.snapDir(), id+TimeFileSuffix) }

func (d *DB) Create(
	version raft.SnapshotVersion,
	index, term uint64,
	configuration raft.Configuration,
	configurationIndex uint64,
	_ raft.Transport,
) (raft.SnapshotSink, error) {
	if err := os.MkdirAll(d.snapDir(), 0755); err != nil {
		return nil, err
	}

	id := SnapshotID(index)
	data, err := json.Marshal(&raft.SnapshotMeta{
		Version:            version,
		ID:                 id,
		Index:              index,
		Term:               term,
		Configuration:      configuration,
		ConfigurationIndex: configurationIndex,
	})
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(d.metaFilePath(id), data, 0644); err != nil {
		return nil, err
	}

	return &fsmDBSink{db: d, id: id}, nil
}

func (d *DB) List() (ret []*raft.SnapshotMeta, err error) {
	entries, err := os.ReadDir(d.snapDir())
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, MetaFileSuffix) {
			continue
		}

		var data []byte
		data, err = os.ReadFile(filepath.Join(d.snapDir(), name))
		if err != nil {
			return nil, err
		}

		var meta raft.SnapshotMeta
		if err = json.Unmarshal(data, &meta); err != nil {
			return
		}

		ret = append(ret, &meta)
	}

	slices.SortStableFunc(ret, func(a, b *raft.SnapshotMeta) int { return strings.Compare(b.ID, a.ID) })
	return
}

func (d *DB) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	var meta raft.SnapshotMeta
	{
		data, err := os.ReadFile(d.metaFilePath(id))
		if err != nil {
			return nil, nil, err
		}
		if err := json.Unmarshal(data, &meta); err != nil {
			return nil, nil, err
		}
	}

	file, err := os.Open(d.snapFilePath(id))
	if err != nil {
		return nil, nil, err
	}

	return &meta, file, nil
}

type BackupDumper interface {
	Backup() error
}

var _ = (BackupDumper)((*fsmDBSink)(nil))

type fsmDBSink struct {
	db *DB
	id string

	snapFile *os.File
	ctx      context.Context
	cancel   context.CancelFunc
}

func (s *fsmDBSink) Backup() error {
	f, err := os.Create(s.db.snapFilePath(s.id))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.snapFile = f
	s.ctx = ctx
	s.cancel = cancel

	ch := make(chan error)
	go func() {
		ch <- s.doBackup()
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *fsmDBSink) doBackup() error {
	maxVersion, err := s.db.Backup(s.snapFile, 0)
	if err != nil {
		return err
	}
	return os.WriteFile(s.db.timeFilePath(s.id), []byte(fmt.Sprint(maxVersion)), 0644)
}

func (*fsmDBSink) Write([]byte) (int, error) { panic("unreachable") }

func (s *fsmDBSink) Close() error { return s.Cancel() }

func (s *fsmDBSink) ID() string { return s.id }

func (s *fsmDBSink) Cancel() error {
	s.cancel()
	return s.snapFile.Close()
}
