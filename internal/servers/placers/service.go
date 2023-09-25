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
	it, err := s.db.NewPrefixIter(DBGroupPrefix)
	if err != nil {
		return nil, err
	}
	for it.First(); it.Valid(); it.Next() {
		ret = append(ret, api.GroupID(it.Key()))
	}
	return &ret, nil
}

func (s *placerServer) GenerateGroup() (*api.GroupID, error) {
	n, err := s.db.Inc(DBGroupCounter)
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
