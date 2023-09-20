package placers

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/logs"
)

const (
	DefaultPlacerDataDir        = "placer-data"
	DefaultPlacerLogCacheCap    = 128
	DefaultPlacerPeersConnPool  = 10
	DefaultPlacerPeersIOTimeout = 15 * time.Second

	PlacerLogDBFile          = "placer.log.bolt"
	PlacerStableDBFile       = "placer.stable.bolt"
	PlacerSnapshotMetaDBFile = "placer.snapmeta.bolt"
)

type PlacerConfig struct {
	Host           string
	ID             raft.ServerID
	PeerMap        *api.AddressMap
	DataDir        string
	LogCacheCap    int
	PeersConnPool  int
	PeersIOTimeout time.Duration

	rawPeers string
}

func DefaultPlacerConfig() *PlacerConfig {
	return &PlacerConfig{
		Host:           api.DefaultPlacerHost,
		ID:             api.DefaultPlacerID,
		PeerMap:        api.DefaultPlacerPeerMap,
		DataDir:        DefaultPlacerDataDir,
		LogCacheCap:    DefaultPlacerLogCacheCap,
		PeersConnPool:  DefaultPlacerPeersConnPool,
		PeersIOTimeout: DefaultPlacerPeersIOTimeout,
	}
}

func (c *PlacerConfig) hcLogger(name string) hclog.Logger {
	return logs.HcLogger(fmt.Sprintf("%s-%s", c.ID, name))
}