package heartbeatleader

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"github.com/codeskyblue/heartbeat"
	"net/http"
	"strconv"
	"syscall"
	"time"
)

const (
	SignalToStopHeatBeatSending = syscall.SIGUSR1
	SignalToTellOfflineNode     = syscall.SIGUSR2
	heartBeatPort               = 30000
	heartBeatingFrequency       = 600 * time.Millisecond
	heartBeatingTimeOut         = 3 * time.Second
)

func OnDisconnect(startNewElection chan bool, server *http.Server) {

	fmt.Println("Leader node is OFFLINE")

	//err := (*server).Close()
	//utility.CheckError(err)

	startNewElection <- true
}

func Receive(nodeId uint, startNewElection chan bool, waitEndElection chan bool) {

	heartBeatServer := heartbeat.NewServer("my-secret", heartBeatingTimeOut)

	http.Handle(fmt.Sprintf("/LeaderHeartBeatTo%dNode", nodeId), heartBeatServer)

	timerFirstHeartBeat := time.NewTimer(heartBeatingTimeOut)

	heartBeatServer.OnConnect = func(identifier string, r *http.Request) {

		timerFirstHeartBeat.Stop()
		fmt.Println(identifier, "Leader is online")
	}

	heartBeatServer.OnDisconnect = func(identifier string) {

		fmt.Println("Leader node is OFFLINE")
		startNewElection <- true
	}

	go func() {
		<-timerFirstHeartBeat.C
		startNewElection <- true
	}()

	http.ListenAndServe(":30000", nil)
}

func Send(senderNodeId uint, recipientNodeId uint, startNewElection chan bool, waitEndElection chan bool) {

	var leaderId uint
	var err error

	senderNodeStringIdentifier := strconv.Itoa(int(senderNodeId))

	client := &heartbeat.Client{
		ServerAddr: fmt.Sprintf("localhost:%d/LeaderHeartBeatTo%dNode", heartBeatPort, recipientNodeId),
		Secret:     "my-secret",
		Identifier: senderNodeStringIdentifier,
	}

	for {

		leaderId, err = zookeeper.GetCurrentLeaderId()
		utility.CheckError(err)

		if leaderId == senderNodeId {
			cancel := client.Beat(heartBeatingFrequency)

			<-waitEndElection

			cancel()
		}

		<-waitEndElection
	}
}

func LeaderElection(nodeId uint, startNewElection chan bool, waitEndElection chan bool) {

	for {
		fmt.Println("Wait for leader election!")
		<-startNewElection

		fmt.Println("Start election")
		actualLeaderId := zookeeper.StartLeaderElection(nodeId)
		fmt.Printf("My ID %d -> Actual Leader ID %d \n", nodeId, actualLeaderId)

		waitEndElection <- true
	}
}
