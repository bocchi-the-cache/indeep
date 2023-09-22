package api

import (
	"errors"

	"github.com/hashicorp/raft"
)

const RpcMemberAskLeaderID RpcID = "ask-leader"

var ErrNotLeader = errors.New("not leader")

type Member interface {
	AskLeaderID() (*raft.ServerID, error)
}

type LeaderChecker interface {
	CheckLeader() error
}
