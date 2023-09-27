package placers

import (
	"fmt"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
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
	it := s.db.Iter([]byte(DBGroupPrefix), false)
	for it.Rewind(); it.Valid(); it.Next() {
		ret = append(ret, api.GroupID(it.Item().KeyCopy(nil)))
	}
	return &ret, nil
}

func (s *placerServer) GenerateGroup() (*api.GroupID, error) {
	n, err := s.db.Inc([]byte(DBGroupCounter))
	if err != nil {
		return nil, err
	}
	id := api.GroupID(fmt.Sprint(DBGroupPrefix, n))
	// FIXME: Generate the group info too.
	return &id, nil
}

func (s *placerServer) CheckLeader() error {
	if _, leaderID := s.rn.LeaderWithID(); s.config.ID != leaderID {
		return fmt.Errorf("%w: %s", api.ErrNotLeader, s.config.ID)
	}
	return nil
}
