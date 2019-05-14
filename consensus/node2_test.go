package consensus

import (
	"SDCC-Project-WorkerNode/consensus/raft/raft"
	"testing"
)

func Test_node3(t *testing.T) {
	raft.Start("2", "127.0.0.1:12002", "127.0.0.1:10000", "127.0.0.1:10002", false)

}
