package placers

import (
	"fmt"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/fsmdb"
)

const (
	DBGroupPrefix  = "g"
	DBGroupCounter = "counterGroup"
)

var _ = (api.Placer)((*placerServer)(nil))

func (s *placerServer) AskLeaderID() (*raft.ServerID, error) {
	_, id := s.rn.LeaderWithID()
	return &id, nil
}

func (s *placerServer) ListGroups() (*[]api.GroupID, error) {
	var ret []api.GroupID
	if err := s.db.View(func(tx *fsmdb.Tx) error {
		it := tx.Iter([]byte(DBGroupPrefix), false)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			ret = append(ret, api.GroupID(it.Item().KeyCopy(nil)))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (s *placerServer) GenerateGroup() (*api.GroupID, error) {
	var id api.GroupID
	if err := s.db.Update(func(tx *fsmdb.Tx) error {
		n, err := tx.Inc([]byte(DBGroupCounter))
		if err != nil {
			return err
		}
		id = api.GroupID(fmt.Sprint(DBGroupPrefix, n))
		// TODO: group range information
		return tx.Set([]byte(id), []byte("TODO"))
	}); err != nil {
		return nil, err
	}
	return &id, nil
}

func (s *placerServer) CheckLeader() error {
	if _, leaderID := s.rn.LeaderWithID(); s.config.ID != leaderID {
		return fmt.Errorf("%w: %s", api.ErrNotLeader, s.config.ID)
	}
	return nil
}
