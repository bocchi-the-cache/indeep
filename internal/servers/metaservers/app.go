package metaservers

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
	"github.com/bocchi-the-cache/indeep/internal/logs"
	"github.com/bocchi-the-cache/indeep/internal/peers"
)

type metaserver struct {
	config *MetaserverConfig
	peers  api.Peers
	server *http.Server

	placerCl api.Placer

	mux api.StreamLayerMux
}

func NewMetaserver(c *MetaserverConfig) api.Server { return &metaserver{config: c} }
func Metaserver() api.Server                       { return NewMetaserver(DefaultMetaserverConfig()) }

func (*metaserver) Name() string { return "metaserver" }

func (m *metaserver) DefineFlags(f *flag.FlagSet) {
	f.StringVar(&m.config.Host, "host", api.DefaultMetaserverHost, "listen host")
	f.StringVar((*string)(&m.config.ID), "id", api.DefaultMetaserverID, "placer ID")
	f.StringVar(&m.config.rawPeers, "peers", api.DefaultMetaserverMultipeerMap.String(), "metaserver peers URL")
	f.StringVar(&m.config.DataDir, "data-dir", DefaultMetaserverDataDir, "data directory")
	f.IntVar(&m.config.LogCacheCap, "logcache-cap", DefaultMetaserverLogCacheCap, "Raft log cache capacity")
	f.IntVar(&m.config.PeersConnPool, "conn-pool", DefaultMetaserverPeersConnPool, "peer connections to pool")
	f.DurationVar(&m.config.PeersIOTimeout, "io-timeout", DefaultMetaserverPeersIOTimeout, "peer IO timeout")
	f.StringVar(&m.config.rawPlacerHosts, "placer-hosts", api.DefaultPlacerHostMap.String(), "placer hosts URL")
}

func (m *metaserver) Setup() error {
	if m.config.rawPeers != "" {
		ps, err := api.ParseAddressMap(m.config.rawPeers)
		if err != nil {
			return err
		}
		m.config.PeerMap = ps
	}
	m.peers = peers.NewPeers(m.config.PeerMap)

	p, err := m.peers.Lookup(m.config.ID)
	if err != nil {
		return err
	}
	// TODO: ServeMux would use this.
	_ = p

	if m.config.rawPlacerHosts != "" {
		ps, err := api.ParseAddressMap(m.config.rawPlacerHosts)
		if err != nil {
			return err
		}
		m.config.Placer.HostMap = ps
	}
	placerCl, err := clients.NewPlacer(&m.config.Placer)
	if err != nil {
		return err
	}
	m.placerCl = placerCl

	groups, err := m.placerCl.ListGroups()
	if err != nil {
		return err
	}
	// TODO: Initialize Raft groups.
	_ = groups

	// TODO
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	m.server = &http.Server{
		Addr:     m.config.Host,
		Handler:  mux,
		ErrorLog: logs.E,
	}

	return nil
}

func (m *metaserver) logDBPath(groupID api.GroupID) string {
	return filepath.Join(m.config.DataDir, fmt.Sprintf("metaserver.log.%s.bolt", groupID))
}

func (m *metaserver) stableDBPath(groupID api.GroupID) string {
	return filepath.Join(m.config.DataDir, fmt.Sprintf("metaserver.stable.%s.bolt", groupID))
}
