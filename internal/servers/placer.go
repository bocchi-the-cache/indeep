package servers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

var ErrPlacerUnknownID = errors.New("unknown place ID")

type PlacerConfig struct {
	ID    api.PeerID
	Peers api.Peers
}

type PlacerServer struct {
	c   *PlacerConfig
	h   *http.Server
	fsm *placerFSM
}

func NewPlacerServer(c *PlacerConfig) (*PlacerServer, error) {
	p := c.Peers.Lookup(c.ID)
	if p == nil {
		return nil, fmt.Errorf("%w: id=%s", ErrPlacerUnknownID, c.ID)
	}

	s := &PlacerServer{
		c: c,
		h: &http.Server{Addr: p.URL().Host},
		fsm: &placerFSM{
			peers:    c.Peers,
			self:     p,
			leader:   p,
			isLeader: true,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(peers.OperationGetMembers, s.Members)
	mux.HandleFunc(peers.OperationAskLeader, s.Leader)
	mux.HandleFunc(peers.OperationLookupMetaService, s.LookupMetaService)
	mux.HandleFunc(peers.OperationAddMetaService, s.AddMetaService)
	mux.HandleFunc(peers.OperationLookupDataService, s.LookupDataService)
	mux.HandleFunc(peers.OperationAddDataService, s.AddDataService)
	s.h.Handler = mux

	return s, nil
}

func (p *PlacerServer) Run() error { return p.h.ListenAndServe() }

func (p *PlacerServer) Shutdown(ctx context.Context) error { return p.h.Shutdown(ctx) }

func (p *PlacerServer) Members(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	jsonhttp.W(w).OK(p.fsm.Members())
}

func (p *PlacerServer) Leader(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	jw := jsonhttp.W(w)
	l, err := p.fsm.Leader(nil)
	if err != nil {
		jw.Err(err)
		return
	}
	jw.OK(l)
}

func (p *PlacerServer) LookupMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = p.fsm.LookupMetaService(nil)
}

func (p *PlacerServer) AddMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = p.fsm.AddMetaService()
}

func (p *PlacerServer) LookupDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = p.fsm.LookupDataService(nil)
}

func (p *PlacerServer) AddDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = p.fsm.AddDataService()
}

type placerFSM struct {
	// FIXME: Any locks here?
	peers    api.Peers
	self     api.Peer
	leader   api.Peer
	isLeader bool
}

func (p *placerFSM) Members() []api.Peer {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) Leader(e api.Peer) (api.Peer, error) {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) LookupMetaService(key api.MetaKey) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) AddMetaService() error {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) LookupDataService(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) AddDataService() error {
	//TODO implement me
	panic("implement me")
}
