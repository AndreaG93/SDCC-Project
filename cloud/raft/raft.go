package raft

import (
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type Config struct {
	Bind    string `json:bind`
	DataDir string `json:data_dir`
}

type Word struct {
	words string
}

func (*Word) Apply(l *raft.Log) interface{} {
	return nil
}

func (*Word) Snapshot() (raft.FSMSnapshot, error) {
	return new(WordSnapshot), nil
}

func (*Word) Restore(snap io.ReadCloser) error {
	return nil
}

type WordSnapshot struct {
	words string
}

func (snap *WordSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (snap *WordSnapshot) Release() {

}

func newRaftCluster(logWriter io.Writer, namePrefix string, n uint, transportHooks TransportHooks) *cluster {

	res := make([]*raftNode, 0, n)
	names := make([]string, 0, n)
	for i := uint(0); i < n; i++ {
		names = append(names, nodeName(namePrefix, i))
	}
	l := log.New(logWriter, "", log.Lmicroseconds)
	transports := newTransports(l)
	for _, i := range names {
		r, err := newRaftNode(log.New(logWriter, i+":", log.Lmicroseconds), transports, transportHooks, names, i)
		if err != nil {
			t.Fatalf("Unable to create raftNode:%v : %v", i, err)
		}
		res = append(res, r)
	}
	return &cluster{
		nodes:        res,
		removedNodes: make([]*raftNode, 0, n),
		applied:      make([]appliedItem, 0, 1024),
		log:          &LoggerAdapter{l},
		transports:   transports,
		hooks:        transportHooks,
	}
}

func CreateRaftNode() (*raft.Raft, error) {

	var raftSnapshotStore *raft.FileSnapshotStore
	var raftBoltDB *raftboltdb.BoltStore
	var raftNetworkTransport *raft.NetworkTransport
	var raftNode *raft.Raft
	var err error

	// cfg.EnableSingleNode = true
	fsm := new(Word)
	fsm.words = "hahaha"

	if raftSnapshotStore, err = raft.NewFileSnapshotStore("./", 1, os.Stdout); err != nil {
		return nil, err
	}

	if raftBoltDB, err = raftboltdb.NewBoltStore(filepath.Join("./", "store.bolt")); err != nil {
		return nil, err
	}

	if raftNetworkTransport, err = raft.NewTCPTransport("localhost:12345", nil, 3, 5*time.Second, os.Stdout); err != nil {
		return nil, err
	}

	if raftNode, err = raft.NewRaft(raft.DefaultConfig(), fsm, raftBoltDB, raftBoltDB, raftSnapshotStore, raftNetworkTransport); err != nil {
		return nil, err
	}

	return raftNode, nil
}
