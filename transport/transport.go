package transport

import (
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type conn struct {
	client *grpc.ClientConn
}

type BadgerGRPC struct {
	address          raft.ServerAddress
	rpcChan          chan raft.RPC
	heartbeatFunc    func(raft.RPC)
	mtx              sync.Mutex
	heartbeatTimeout time.Duration

	connections map[raft.ServerID]*conn
}

func New(addr raft.ServerAddress) BadgerGRPC {

	return BadgerGRPC{
		address:     addr,
		rpcChan:     make(chan raft.RPC),
		connections: map[raft.ServerID]*conn{},
	}

}

func (b *BadgerGRPC) Transport() Api {
	return Api{badgerGRPC: b}

}
