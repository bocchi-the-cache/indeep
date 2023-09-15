package placers

import (
	"context"
	"errors"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
)

func (s *placerServer) ListenAndServe() error { return s.server.ListenAndServe() }
func (s *placerServer) Shutdown(ctx context.Context) error {
	return errors.Join(s.rn.Shutdown().Error(), s.server.Shutdown(ctx))
}

func (s *placerServer) HandleGetMembers() (raft.Configuration, error) {
	return s.GetMembers().Configuration(), nil
}

func (s *placerServer) HandleAskLeader() (api.Peer, error) {
	return s.AskLeader(nil)
}
