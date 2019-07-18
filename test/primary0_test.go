package test

import (
	"SDCC-Project-WorkerNode/consensus/raft/raft"
	"sync"
	"testing"
)

func Test_primary0(t *testing.T) {

	var myWaitGroup sync.WaitGroup

	myWaitGroup.Add(1)

	go raft.Start("0", "127.0.0.1:12000", "", "127.0.0.1:10000", true)
	//raft.Start("1", "127.0.0.1:12001", "127.0.0.1:10000", "127.0.0.1:10001", false)

	myWaitGroup.Wait()
}
