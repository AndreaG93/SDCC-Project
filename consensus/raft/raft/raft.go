package raft

import (
	"SDCC-Project-WorkerNode/consensus/http"
	"SDCC-Project-WorkerNode/consensus/raft/raftnode"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/utility"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	raftDirectory = "./raftDirectory"
)

func Start(nodeId string, raftBindAddress string, joinAddress string, httpAddress string, bootstrapCluster bool) {

	if bootstrapCluster {
		zookeeper.Init()
	}

	var err error

	err = os.MkdirAll(raftDirectory, 0700)
	utility.CheckError(err)

	raftNode := raftnode.New(raftBindAddress, raftDirectory)
	raftNode.Start(nodeId, bootstrapCluster)

	h := httpd.New(httpAddress, raftNode)
	if err := h.Start(); err != nil {
		log.Fatalf("failed to start HTTP service: %raftLogStore", err.Error())
	}

	// If join was specified, make the join request.
	if joinAddress != "" {

		fmt.Println("Node joining")

		if err := join(joinAddress, raftBindAddress, nodeId); err != nil {
			log.Fatalf("failed to join node at %raftLogStore: %raftLogStore", joinAddress, err.Error())
		}
	}

	go func() {

		for {

			if raftNode.IsLeader() {

				zookeeper.SetActualClusterLeaderAddress(raftNode.Leader())
				fmt.Printf("Actual Leader ID %s \n", nodeId)
			}

			<-raftNode.LeaderCh()
		}
	}()

	log.Println("RAFT CONSENSUS algorithm is working!")
	waitInterrupt()
}

func waitInterrupt() {

	terminateNodeWork := make(chan os.Signal, 1)
	signal.Notify(terminateNodeWork, os.Interrupt)
	<-terminateNodeWork
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
