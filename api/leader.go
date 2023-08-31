package api

type Leadership interface {
	IsLeader(e Endpoint) (bool, error)
}
