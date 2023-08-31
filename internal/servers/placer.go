package servers

import (
	"context"
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/endpoints"
)

type PlacerConfig struct {
	Addr string
}

type PlacerServer struct {
	c *PlacerConfig
	h *http.Server
	p api.Placer
}

func NewPlacerServer(c *PlacerConfig) *PlacerServer {
	p := &PlacerServer{c: c, h: &http.Server{Addr: c.Addr}, p: new(placerFSM)}
	mux := http.NewServeMux()
	mux.HandleFunc(endpoints.OperationGetMembers, p.Members)
	mux.HandleFunc(endpoints.OperationAskLeader, p.Leader)
	mux.HandleFunc(endpoints.OperationLookupMetaService, p.LookupMetaService)
	mux.HandleFunc(endpoints.OperationAddMetaService, p.AddMetaService)
	mux.HandleFunc(endpoints.OperationLookupDataService, p.LookupDataService)
	mux.HandleFunc(endpoints.OperationAddDataService, p.AddDataService)
	p.h.Handler = mux
	return p
}

func (p *PlacerServer) Shutdown(ctx context.Context) error { return p.h.Shutdown(ctx) }

func (p *PlacerServer) Members(w http.ResponseWriter, r *http.Request) {
	// TODO
	p.p.Members()
}

func (p *PlacerServer) Leader(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = p.p.Leader(nil)
}

func (p *PlacerServer) LookupMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = p.p.LookupMetaService(nil)
}

func (p *PlacerServer) AddMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = p.p.AddMetaService()
}

func (p *PlacerServer) LookupDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = p.p.LookupDataService(nil)
}

func (p *PlacerServer) AddDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = p.p.AddDataService()
}

type placerFSM struct{}

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
