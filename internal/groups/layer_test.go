package groups

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/logs"
)

func newStreamLayerMux(t *testing.T, local api.NodeHost) api.StreamLayerMux {
	m, err := NewStreamLayerMux(local)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func newBoltDB(t *testing.T, name string, id raft.ServerID) *raftboltdb.BoltStore {
	id = raft.ServerID(strings.ReplaceAll(string(id), "/", "_"))
	db, err := raftboltdb.New(raftboltdb.Options{Path: fmt.Sprintf("%s.%s.bolt", name, id)})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func bootstrapNode(t *testing.T, s raft.StreamLayer, c raft.Configuration, id raft.ServerID) {
	config := raft.DefaultConfig()
	config.LocalID = id
	snaps, err := raft.NewFileSnapshotStoreWithLogger(".", 3, logs.HashiCorp("snaps"))
	if err != nil {
		t.Fatal(err)
	}
	trans := raft.NewNetworkTransportWithLogger(s, 10, 15*time.Second, logs.HashiCorp("trans"))
	r, err := raft.NewRaft(
		config,
		new(raft.MockFSM),
		newBoltDB(t, "log", id),
		newBoltDB(t, "stable", id),
		snaps,
		trans,
	)
	if err != nil {
		t.Fatal(err)
	}
	r.BootstrapCluster(c)
	go func() {
		time.Sleep(10 * time.Second)
		leader, leaderID := r.LeaderWithID()
		logs.S.Info("Raft group info", "leader", leader, "leaderID", leaderID, "c", r.GetConfiguration().Configuration())
	}()
}

func TestNewTransMux(t *testing.T) {
	c1 := raft.Configuration{Servers: []raft.Server{
		{ID: "n1/g1", Address: "127.0.0.1:10001/g1"},
		{ID: "n2/g1", Address: "127.0.0.1:10002/g1"},
		{ID: "n3/g1", Address: "127.0.0.1:10003/g1"},
	}}
	n1 := newStreamLayerMux(t, "127.0.0.1:10001")
	n1g1 := n1.NetworkLayer("g1")
	n1g2 := n1.NetworkLayer("g2")
	n1g3 := n1.NetworkLayer("g3")

	c2 := raft.Configuration{Servers: []raft.Server{
		{ID: "n1/g2", Address: "127.0.0.1:10001/g2"},
		{ID: "n2/g2", Address: "127.0.0.1:10002/g2"},
		{ID: "n3/g2", Address: "127.0.0.1:10003/g2"},
	}}
	n2 := newStreamLayerMux(t, "127.0.0.1:10002")
	n2g1 := n2.NetworkLayer("g1")
	n2g2 := n2.NetworkLayer("g2")
	n2g3 := n2.NetworkLayer("g3")

	c3 := raft.Configuration{Servers: []raft.Server{
		{ID: "n1/g3", Address: "127.0.0.1:10001/g3"},
		{ID: "n2/g3", Address: "127.0.0.1:10002/g3"},
		{ID: "n3/g3", Address: "127.0.0.1:10003/g3"},
	}}
	n3 := newStreamLayerMux(t, "127.0.0.1:10003")
	n3g1 := n3.NetworkLayer("g1")
	n3g2 := n3.NetworkLayer("g2")
	n3g3 := n3.NetworkLayer("g3")

	bootstrapNode(t, n1g1, c1, c1.Servers[0].ID)
	bootstrapNode(t, n1g2, c1, c2.Servers[0].ID)
	bootstrapNode(t, n1g3, c1, c3.Servers[0].ID)

	bootstrapNode(t, n2g1, c2, c1.Servers[1].ID)
	bootstrapNode(t, n2g2, c2, c2.Servers[1].ID)
	bootstrapNode(t, n2g3, c2, c3.Servers[1].ID)

	bootstrapNode(t, n3g1, c3, c1.Servers[2].ID)
	bootstrapNode(t, n3g2, c3, c2.Servers[2].ID)
	bootstrapNode(t, n3g3, c3, c3.Servers[2].ID)

	select {}
}
