package servers

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/jsonhttp"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

const (
	DefaultPlacerID   = "placer0"
	DefaultPlacerHost = "127.0.0.1:11402"
)

var (
	ErrPlacerUnknownID = errors.New("unknown placer ID")

	DefaultPlacerRawPeers = (&url.URL{
		Scheme: peers.DefaultScheme,
		Host:   DefaultPlacerHost,
		Path:   peers.IDsPrefix + DefaultPlacerID,
	}).String()
)

type PlacerConfig struct {
	ID    api.PeerID
	Peers api.Peers

	rawPeers string
}

type placerServer struct {
	c   *PlacerConfig
	h   *http.Server
	fsm *placerFSM
}

func NewPlacer(c *PlacerConfig) api.App { return &placerServer{c: c} }
func Placer() api.App                   { return NewPlacer(new(PlacerConfig)) }

func (*placerServer) Name() string { return "placer" }

func (s *placerServer) DefineFlags(f *flag.FlagSet) {
	f.StringVar((*string)(&s.c.ID), "id", DefaultPlacerID, "placer ID")
	f.StringVar(&s.c.rawPeers, "peers", DefaultPlacerRawPeers, "full placer peers")
}

func (s *placerServer) Initialize() error {
	ps, err := peers.ParsePeers(s.c.rawPeers)
	if err != nil {
		return err
	}
	s.c.Peers = ps

	p := ps.Lookup(s.c.ID)
	if p == nil {
		return fmt.Errorf("%w: id=%s", ErrPlacerUnknownID, s.c.ID)
	}

	// TODO
	s.fsm = &placerFSM{
		peers:    ps,
		self:     p,
		leader:   p,
		isLeader: true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(peers.OperationGetMembers, s.Members)
	mux.HandleFunc(peers.OperationAskLeader, s.Leader)
	mux.HandleFunc(peers.OperationLookupMetaService, s.LookupMetaService)
	mux.HandleFunc(peers.OperationAddMetaService, s.AddMetaService)
	mux.HandleFunc(peers.OperationLookupDataService, s.LookupDataService)
	mux.HandleFunc(peers.OperationAddDataService, s.AddDataService)
	s.h = &http.Server{Addr: p.URL().Host, Handler: mux}

	return nil
}

func (s *placerServer) Run() error                         { return s.h.ListenAndServe() }
func (s *placerServer) Shutdown(ctx context.Context) error { return s.h.Shutdown(ctx) }

func (s *placerServer) Members(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	jsonhttp.W(w).OK(s.fsm.Members())
}

func (s *placerServer) Leader(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	jw := jsonhttp.W(w)
	l, err := s.fsm.Leader(nil)
	if err != nil {
		jw.Err(err)
		return
	}
	jw.OK(l)
}

func (s *placerServer) LookupMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = s.fsm.LookupMetaService(nil)
}

func (s *placerServer) AddMetaService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = s.fsm.AddMetaService()
}

func (s *placerServer) LookupDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_, _ = s.fsm.LookupDataService(nil)
}

func (s *placerServer) AddDataService(w http.ResponseWriter, r *http.Request) {
	// TODO
	_ = s.fsm.AddDataService()
}

type placerFSM struct {
	// FIXME: Any locks here?
	peers    api.Peers
	self     api.Peer
	leader   api.Peer
	isLeader bool
}

func (p *placerFSM) Members() api.Peers {
	return p.peers
}

func (p *placerFSM) Leader(api.Peer) (api.Peer, error) {
	return p.leader, nil
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
