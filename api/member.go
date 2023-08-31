package api

type LeadershipRequest struct{}

type Membership interface {
	Members() []Endpoint
	Leader() Endpoint
}
