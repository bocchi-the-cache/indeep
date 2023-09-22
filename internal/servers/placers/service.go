package placers

import (
	"fmt"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/utils"
)

const DBGroupPrefix = "G"

var _ = (api.Placer)((*placerServer)(nil))

func (s *placerServer) AskLeaderID() (*raft.ServerID, error) {
	_, id := s.rn.LeaderWithID()
	return &id, nil
}

func (s *placerServer) ListGroups() (*[]api.GroupID, error) {
	var ret []api.GroupID
	it, err := s.db.NewIter(utils.PrefixIterOptions(DBGroupPrefix))
	if err != nil {
		return nil, err
	}
	for it.First(); it.Valid(); it.Next() {
		ret = append(ret, api.GroupID(it.Key()[len(DBGroupPrefix):]))
	}
	return &ret, nil
}

func (s *placerServer) GenerateGroup() (*api.GroupID, error) {
	// TODO
	panic("TODO")
}

func (s *placerServer) CheckLeader() error {
	if _, leaderID := s.rn.LeaderWithID(); s.config.ID != leaderID {
		return fmt.Errorf("%w: %s", api.ErrNotLeader, s.config.ID)
	}
	return nil
}
