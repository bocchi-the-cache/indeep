package api

import "errors"

var (
	ErrEmptyMembers = errors.New("empty members")
)

type Member interface {
	Members() Peers
	Leader(e Peer) (Peer, error)
}

func AskLeader(m Member) (Peer, error) {
	for _, e := range m.Members().Peers() {
		return m.Leader(e)
	}
	return nil, ErrEmptyMembers
}
