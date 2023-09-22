package api

import (
	"errors"

	"github.com/hashicorp/raft"
)

const (
	RpcMemberGetMembers RpcID = "get-members"
	RpcMemberAskLeader  RpcID = "ask-leader"
)

var ErrEmptyMembers = errors.New("empty members")

type Member interface {
	Peers() Peers
	AskLeaderID(e Peer) (raft.ServerID, error)
}

func AskLeaderID(m Member) (raft.ServerID, error) {
	for _, e := range m.Peers().Peers() {
		return m.AskLeaderID(e)
	}
	return "", ErrEmptyMembers
}
