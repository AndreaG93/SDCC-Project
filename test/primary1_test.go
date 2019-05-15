package test

import (
	"SDCC-Project-WorkerNode/consensus/raft/raft"
	"testing"
)

func Test_primary1(t *testing.T) {
	raft.Start("1", "127.0.0.1:12001", "127.0.0.1:10000", "127.0.0.1:10001", false)
}
