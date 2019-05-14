package consensus

import (
	"SDCC-Project-WorkerNode/consensus/raft/raft"
	"testing"
)

func Test_node0(t *testing.T) {
	raft.Start("0", "127.0.0.1:12000", "", "127.0.0.1:10000", true)
}
