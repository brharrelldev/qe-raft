package main

import (
	"github.com/hashicorp/raft"
	"log"
	"os"
	"path/filepath"
)

func main() {

	conf := &raft.Config{}

	inv := &qeInv{}

	snapShotDir := filepath.Join("snapshot", "id")

	stableStore, err := NewBadgerStore("stable")
	if err != nil {
		log.Fatal(err)
	}

	logstore, err := NewBadgerStore("logs")
	if err != nil {
		log.Fatal(err)
	}

	snap, err := raft.NewFileSnapshotStore(snapShotDir, 3, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	cluster, err := raft.NewRaft(conf, inv, logstore, stableStore, snap, nil)
	if err != nil {
		log.Fatal(err)
	}

	fut := cluster.BootstrapCluster(raft.Configuration{Servers: []raft.Server{
		{
			Suffrage: 1,
			ID:       "id",
			Address:  "127.0.0.1",
		},
	}})

	if err := fut.Error(); err != nil {
		log.Fatal(err)
	}

}
