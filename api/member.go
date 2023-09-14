package api

import (
	"errors"
)

var (
	ErrEmptyMembers = errors.New("empty members")
)

type Member interface {
	GetMembers() Peers
	AskLeader(e Peer) (Peer, error)
}

func AskLeader(m Member) (Peer, error) {
	for _, e := range m.GetMembers().Peers() {
		return m.AskLeader(e)
	}
	return nil, ErrEmptyMembers
}
