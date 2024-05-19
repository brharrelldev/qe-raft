package main

import (
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	conf := &raft.Config{
		ProtocolVersion:    raft.ProtocolVersionMax,
		LocalID:            "test1",
		HeartbeatTimeout:   time.Millisecond * 5000,
		ElectionTimeout:    time.Second * 10,
		CommitTimeout:      time.Second * 1,
		MaxAppendEntries:   100,
		SnapshotInterval:   time.Millisecond * 500,
		LeaderLeaseTimeout: time.Second * 5,
	}

	inv := &qeInv{}

	snapShotDir := filepath.Join("snapshot", "id")

	stableStore, err := NewBadgerStore("stable")
	if err != nil {
		logger.Fatal("Failed to open stable store", zap.Error(err))
	}

	logstore, err := NewBadgerStore("logs")
	if err != nil {
		logger.Fatal("Failed to open logs store", zap.Error(err))
	}

	snap, err := raft.NewFileSnapshotStore(snapShotDir, 3, os.Stderr)
	if err != nil {
		logger.Fatal("Failed to create snapshot store", zap.Error(err))
	}

	creds := insecure.NewCredentials()

	tp := transport.New("127.0.0.1", []grpc.DialOption{grpc.WithTransportCredentials(creds)})

	cluster, err := raft.NewRaft(conf, inv, logstore, stableStore, snap, tp.Transport())
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
