package main

import (
	"fmt"
	"github.com/hashicorp/raft"
)

var _ raft.FSMSnapshot = snapshot{}

type snapshot struct {
	data []byte
}

func (s snapshot) Persist(sink raft.SnapshotSink) error {
	//TODO implement me
	if _, err := sink.Write(s.data); err != nil {
		return fmt.Errorf("error persisting data %v", err)
	}

	return nil
}

func (s snapshot) Release() {

}
