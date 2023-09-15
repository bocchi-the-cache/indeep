package placers

import "github.com/bocchi-the-cache/indeep/api"

func (s *placerServer) GetMembers() api.Peers { return s.peers }

func (s *placerServer) AskLeader(api.Peer) (api.Peer, error) {
	_, id := s.rn.LeaderWithID()
	return s.peers.Lookup(id), nil
}

func (s *placerServer) LookupMetaService(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) AddMetaService() error {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) LookupDataService(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (s *placerServer) AddDataService() error {
	//TODO implement me
	panic("implement me")
}
