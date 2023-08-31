package servers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/endpoints"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
)

var ErrPlacerUnknownID = errors.New("unknown place ID")

type PlacerConfig struct {
	ID    string
	Peers api.EndpointMap
}

type PlacerServer struct {
	c   *PlacerConfig
	h   *http.Server
	fsm *placerFSM
}

func NewPlacerServer(c *PlacerConfig) (*PlacerServer, error) {
	m := c.Peers.Endpoints()
	e, ok := m[c.ID]
	if !ok {
		return nil, fmt.Errorf("%w: id=%s", ErrPlacerUnknownID, c.ID)
	}

	p := &PlacerServer{
		c: c,
		h: &http.Server{Addr: e.URL().Host},
		fsm: &placerFSM{
			peers:    c.Peers,
			self:     e,
			leader:   e,
			isLeader: true,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(endpoints.OperationGetMembers, p.Members)
	mux.HandleFunc(endpoints.OperationAskLeader, p.Leader)
	mux.HandleFunc(endpoints.OperationLookupMetaService, p.LookupMetaService)
	mux.HandleFunc(endpoints.OperationAddMetaService, p.AddMetaService)
	mux.HandleFunc(endpoints.OperationLookupDataService, p.LookupDataService)
	mux.HandleFunc(endpoints.OperationAddDataService, p.AddDataService)
	p.h.Handler = mux

	return p, nil
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
	peers    api.EndpointMap
	self     api.Endpoint
	leader   api.Endpoint
	isLeader bool
}

func (p *placerFSM) Members() []api.Endpoint {
	//TODO implement me
	panic("implement me")
}

func (p *placerFSM) Leader(e api.Endpoint) (api.Endpoint, error) {
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
