package api

import "errors"

var (
	ErrEmptyMembers = errors.New("empty members")
)

type Member interface {
	Members() []Endpoint
	Leader(e Endpoint) (Endpoint, error)
}

func AskLeader(m Member) (Endpoint, error) {
	if members := m.Members(); len(members) > 0 {
		return m.Leader(members[0])
	}
	return nil, ErrEmptyMembers
}
