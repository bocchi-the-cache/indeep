package metaservers

import (
	"time"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/clients"
)

const (
	DefaultMetaserverDataDir        = "metaserver-data"
	DefaultMetaserverSnapshotRetain = 10
	DefaultMetaserverLogCacheCap    = 128
	DefaultMetaserverPeersConnPool  = 10
	DefaultMetaserverPeersIOTimeout = 15 * time.Second
)

type MetaserverConfig struct {
	Host           string
	ID             raft.ServerID
	PeerMap        *api.AddressMap
	DataDir        string
	SnapshotRetain int
	LogCacheCap    int
	PeersConnPool  int
	PeersIOTimeout time.Duration
	Placer         clients.PlacerConfig

	rawPeers       string
	rawPlacerHosts string
}

func DefaultMetaserverConfig() *MetaserverConfig {
	return &MetaserverConfig{
		Host: api.DefaultMetaserverHost,
		Placer: clients.PlacerConfig{
			HostMap:       api.DefaultPlacerHostMap,
			ClientTimeout: 15 * time.Second,
		},
	}
}
