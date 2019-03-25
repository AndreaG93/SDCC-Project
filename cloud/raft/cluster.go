package raft

import (
	"fmt"
	"github.com/hashicorp/raft"
	"io"
	"log"
	"testing"
	"time"
)

type cluster struct {
	nodes            []*raftNode
	removedNodes     []*raftNode
	lastApplySuccess raft.ApplyFuture
	lastApplyFailure raft.ApplyFuture
	applied          []appliedItem
	log              Logger
	transports       *transports
	hooks            TransportHooks
}

// Logger is abstract type for debug log messages
type Logger interface {
	Log(v ...interface{})
	Logf(s string, v ...interface{})
}

// LoggerAdapter allows a log.Logger to be used with the local Logger interface
type LoggerAdapter struct {
	log *log.Logger
}

// Log a message to the contained debug log
func (a *LoggerAdapter) Log(v ...interface{}) {
	a.log.Print(v...)
}

// Logf will record a formatted message to the contained debug log
func (a *LoggerAdapter) Logf(s string, v ...interface{}) {
	a.log.Printf(s, v...)
}

func newRaftCluster(t *testing.T, logWriter io.Writer, namePrefix string, n uint, transportHooks TransportHooks) *cluster {
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

func (c *cluster) CreateAndAddNode(logWriter io.Writer, namePrefix string, nodeNum uint) error {

	name := nodeName(namePrefix, nodeNum)
	rn, err := newRaftNode(log.New(logWriter, name+":", log.Lmicroseconds), c.transports, c.hooks, nil, name)
	if err != nil {
		fmt.Print("Unable to create raftNode:%v : %v", name, err)
	}
	c.nodes = append(c.nodes, rn)
	f := c.Leader(time.Minute).raft.AddVoter(raft.ServerID(name), raft.ServerAddress(name), 0, 0)

	return f.Error()
}

func nodeName(prefix string, num uint) string {
	return fmt.Sprintf("%v_%d", prefix, num)
}

func (c *cluster) Leader(timeout time.Duration) *raftNode {
	start := time.Now()
	for true {
		for _, n := range c.nodes {
			if n.raft.State() == raft.Leader {
				return n
			}
		}
		if time.Now().Sub(start) > timeout {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return nil
}
