package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
)

type BadgeStore struct {
	db *badger.DB
}

var _ raft.LogStore = &BadgeStore{}
var _ raft.StableStore = &BadgeStore{}

func NewBadgerStore(path string) (*BadgeStore, error) {

	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("error creating badger database %v", err)
	}

	return &BadgeStore{
		db: db,
	}, nil

}

func (b BadgeStore) FirstIndex() (uint64, error) {

	var index uint64
	if err := b.db.View(func(txn *badger.Txn) error {

		prefix := []byte{0x0}

		iter := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
		})

		defer iter.Close()

		iter.Seek(prefix)

		if iter.ValidForPrefix(prefix) {
			index = binary.BigEndian.Uint64(iter.Item().Key()[1:])
		}

		return nil

	}); err != nil {
		return 0, fmt.Errorf("error getting first index from logs %v", err)
	}

	return index, nil
}

func (b BadgeStore) LastIndex() (uint64, error) {

	var index uint64

	prefix := []byte{0x0}

	if err := b.db.View(func(txn *badger.Txn) error {

		iter := txn.NewIterator(badger.IteratorOptions{
			Reverse:        true,
			PrefetchValues: false,
		})

		iter.Seek(append(prefix, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff))

		if iter.ValidForPrefix(prefix) {

			index = binary.BigEndian.Uint64(iter.Item().Key()[1:])

		}

		return nil

	}); err != nil {

		return 0, fmt.Errorf("error getting last index from logs %v", err)

	}

	return index, nil
}

func (b BadgeStore) GetLog(index uint64, log *raft.Log) error {

	prefix := []byte{0x0}

	if err := b.db.View(func(txn *badger.Txn) error {

		buf := make([]byte, 8)

		binary.BigEndian.PutUint64(buf, index)

		data, err := txn.Get(append(prefix, buf...))
		if err != nil {
			return err
		}

		val, err := data.ValueCopy(nil)
		if err != nil {
			return err
		}

		logBuff := bytes.NewBuffer(val)
		msgPack := codec.MsgpackHandle{}

		dec := codec.NewDecoder(logBuff, &msgPack)
		if err := dec.Decode(log); err != nil {
			return fmt.Errorf("error decoding log %v", err)
		}

		return nil

	}); err != nil {
		return fmt.Errorf("error retrieving logs from store %v", err)
	}

	return nil
}

func (b BadgeStore) StoreLog(log *raft.Log) error {

	prefix := []byte{0x0}

	logBuff := bytes.NewBuffer(nil)
	msgPack := codec.MsgpackHandle{}
	enc := codec.NewEncoder(logBuff, &msgPack)

	if err := enc.Encode(log); err != nil {
		return fmt.Errorf("error encodiong message pack %v", err)
	}

	if err := b.db.Update(func(txn *badger.Txn) error {

		buf := make([]byte, 8)

		binary.BigEndian.PutUint64(buf, log.Index)
		if err := txn.Set(append(prefix, buf...), log.Data); err != nil {
			return fmt.Errorf("error appending log %v", err)
		}
		return nil

	}); err != nil {
		return fmt.Errorf("error wrting entry to log %v", err)

	}

	return nil
}

func (b BadgeStore) StoreLogs(logs []*raft.Log) error {

	prefix := []byte{0x0}

	txn := b.db.NewTransaction(true)
	for _, log := range logs {

		buf := bytes.NewBuffer(nil)

		msgPack := codec.MsgpackHandle{}
		enc := codec.NewEncoder(buf, &msgPack)
		if err := enc.Encode(log); err != nil {
			fmt.Println("error encoding log")
			continue
		}

		if err := b.db.Update(func(txn *badger.Txn) error {

			keyBuf := make([]byte, 8)

			binary.BigEndian.PutUint64(keyBuf, log.Index)
			if err := txn.Set(append(prefix, keyBuf...), buf.Bytes()); err != nil {
				return fmt.Errorf("write operation failed %v", err)
			}

			return nil

		}); err != nil {
			return fmt.Errorf("error appending bulk entries to logs %v", err)

		}

	}

	if err := txn.Commit(); err != nil {
		return fmt.Errorf("error committing transaction %v", err)
	}

	return nil
}

func (b BadgeStore) DeleteRange(min, max uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b BadgeStore) Set(key []byte, val []byte) error {
	//TODO implement me
	panic("implement me")
}

func (b BadgeStore) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b BadgeStore) SetUint64(key []byte, val uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b BadgeStore) GetUint64(key []byte) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
