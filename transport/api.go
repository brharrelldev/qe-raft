package transport

import (
	"github.com/hashicorp/raft"
	"io"
)

type Api struct {
	badgerGRPC *BadgerGRPC
}

func (a Api) Consumer() <-chan raft.RPC {
	//TODO implement me
	return a.badgerGRPC.rpcChan
}

func getPeer(id *raft.ServerID, target raft.ServerAddress) {

}

func (a Api) LocalAddr() raft.ServerAddress {
	//TODO implement me
	return a.badgerGRPC.address
}

func (a Api) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	//TODO implement me
	panic("implement me")
}

func (a Api) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	//TODO implement me
	panic("implement me")
}

func (a Api) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	//TODO implement me
	panic("implement me")
}

func (a Api) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (a Api) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	//TODO implement me
	panic("implement me")
}

func (a Api) DecodePeer(bytes []byte) raft.ServerAddress {
	//TODO implement me
	panic("implement me")
}

func (a Api) SetHeartbeatHandler(cb func(rpc raft.RPC)) {
	//TODO implement me
	panic("implement me")
}

func (a Api) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	//TODO implement me
	panic("implement me")
}
