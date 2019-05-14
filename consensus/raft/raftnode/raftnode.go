package raftnode

import (
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type RaftNode struct {
	raft             *raft.Raft
	raftBindAddress  string
	logDatabaseStore *raftboltdb.BoltStore
	raftDirectory    string
}

type fsm RaftNode

func (fsm) Apply(*raft.Log) interface{} {
	panic("implement me")
}

func (fsm) Snapshot() (raft.FSMSnapshot, error) {
	panic("implement me")
}

func (fsm) Restore(io.ReadCloser) error {
	panic("implement me")
}

func New(raftBindAddress string, raftDirectory string) *RaftNode {

	output := new(RaftNode)

	(*output).raftBindAddress = raftBindAddress
	(*output).raftDirectory = raftDirectory

	return output
}

func (obj *RaftNode) Start(nodeId string, bootstrapCluster bool) {

	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeId)

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", (*obj).raftBindAddress)
	utility.CheckError(err)

	transport, err := raft.NewTCPTransport((*obj).raftBindAddress, addr, 3, 10*time.Second, os.Stderr)
	utility.CheckError(err)

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore((*obj).raftDirectory, retainSnapshotCount, os.Stderr)
	utility.CheckError(err)

	// Create the log store and stable store.
	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	// Instantiate the Raft systems.
	(*obj).raft, err = raft.NewRaft(config, (*fsm)(obj), logStore, stableStore, snapshots, transport)

	// Boot consensus cluster.
	if bootstrapCluster {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		(*obj).raft.BootstrapCluster(configuration)
	}

	utility.CheckError(err)
}

func (obj *RaftNode) IsLeader() bool {
	addr := (*obj).raft.Leader()
	if addr == "" {
		return false
	} else {
		return string(addr) == (*obj).raftBindAddress
	}
}

func (obj *RaftNode) JoinTo(nodeID, addr string) {
	fmt.Printf("Received join request for remote node as %s", addr)

	configFuture := (*obj).raft.GetConfiguration()
	utility.CheckError(configFuture.Error())

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				fmt.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
			}

			future := (*obj).raft.RemoveServer(srv.ID, 0, 0)
			utility.CheckError(future.Error())
		}
	}

	indexFuture := (*obj).raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	utility.CheckError(indexFuture.Error())

	fmt.Printf("Node at %s joined successfully", addr)
}

func (obj *RaftNode) LeaderCh() <-chan bool {
	return (*obj).raft.LeaderCh()
}

func (obj *RaftNode) Leader() string {
	return string((*obj).raft.Leader())
}
