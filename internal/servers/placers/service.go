package placers

import (
	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/utils"
)

const DBGroupPrefix = "G"

var _ = (api.Placer)((*placerServer)(nil))

func (s *placerServer) Peers() api.Peers { return s.peers }

func (s *placerServer) AskLeaderID(api.Peer) (raft.ServerID, error) {
	_, id := s.rn.LeaderWithID()
	return id, nil
}

func (s *placerServer) ListGroups() (ret []api.GroupID, err error) {
	it, err := s.db.NewIter(utils.PrefixIterOptions(DBGroupPrefix))
	if err != nil {
		return
	}
	for it.First(); it.Valid(); it.Next() {
		ret = append(ret, api.GroupID(it.Key()[len(DBGroupPrefix):]))
	}
	return
}
