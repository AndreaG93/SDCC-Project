package test

import (
	"SDCC-Project-WorkerNode/consensus/raft/raft"
	"testing"
)

func Test_primary3(t *testing.T) {
	raft.Start("3", "127.0.0.1:12003", "127.0.0.1:10000", "127.0.0.1:10003", false)

}
