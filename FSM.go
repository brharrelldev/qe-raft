package main

import (
	"fmt"
	"github.com/hashicorp/raft"
	"io"
)

var _ raft.FSM = &qeInv{}

type qeInv struct {
	data []byte
}

func (bls *qeInv) Apply(log *raft.Log) interface{} {

	bls.data = log.Data

	return nil
}

func (bls *qeInv) Snapshot() (raft.FSMSnapshot, error) {
	return snapshot{data: bls.data}, nil

}

func (bls *qeInv) Restore(snap io.ReadCloser) error {
	buf := make([]byte, 1024)

	n, err := snap.Read(buf)
	if err != nil {
		return fmt.Errorf("error reading from buffer for restore %v", err)
	}

	bls.data = buf[:n]

	return nil
}
