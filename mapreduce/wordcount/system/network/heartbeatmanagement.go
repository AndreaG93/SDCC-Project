package network

import (
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"github.com/codeskyblue/heartbeat"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	SignalToStopHeatBeatSending = syscall.SIGUSR1
	SignalToTellOfflineNode     = syscall.SIGUSR2
	heartBeatPort               = "30000"
	heartBeatingFrequency       = 400 * time.Millisecond
	heartBeatingTimeOut         = 3 * time.Second
)

func StopHeartBeat() {
	err := syscall.Kill(syscall.Getpid(), SignalToStopHeatBeatSending)
	utility.CheckError(err)
}

func SendHeartbeatToPrimaryBackup(leaderNodeId uint, primaryBackup uint) {
	SendHeartbeatToNode(leaderNodeId, primaryBackup, "/LeaderHeartBeatTo")
}

func SendHeartbeatFromWorkerToLeader(senderNodeId uint, workerNodeId uint) {
	SendHeartbeatToNode(senderNodeId, workerNodeId, "/WorkerHeartBeatToLeader")
}

func SendHeartbeatToNode(senderNodeId uint, recipientNodeId uint, addressPath string) {

	go func() {

		senderNodeStringIdentifier := strconv.Itoa(int(senderNodeId))

		recipientNodeStringIdentifier := strconv.Itoa(int(recipientNodeId))

		client := &heartbeat.Client{
			ServerAddr: "localhost" + ":" + heartBeatPort + addressPath + recipientNodeStringIdentifier,
			Secret:     "my-secret",
			Identifier: senderNodeStringIdentifier,
		}
		cancel := client.Beat(heartBeatingFrequency)

		stopHeartBeatChannel := make(chan os.Signal, 1)
		signal.Notify(stopHeartBeatChannel, SignalToStopHeatBeatSending)

		<-stopHeartBeatChannel

		cancel()
	}()
}

func ReceiveHeartbeatFromLeaderNode(nodeId uint, leaderId uint, alreadyRegister bool) {

	go func() {

		nodeStringIdentifier := strconv.Itoa(int(nodeId))

		timerFirstHeartBeat := time.NewTimer(heartBeatingTimeOut)

		heartBeatServer := heartbeat.NewServer("my-secret", heartBeatingTimeOut) // secret: my-secret, timeout: 15s

		if !alreadyRegister {
			http.Handle("/LeaderHeartBeatTo"+nodeStringIdentifier, heartBeatServer)
		}

		server := &http.Server{Addr: ":" + heartBeatPort, Handler: heartBeatServer}

		heartBeatServer.OnConnect = func(identifier string, r *http.Request) {

			timerFirstHeartBeat.Stop()

			fmt.Println(identifier, "is online")

			id, _ := strconv.Atoi(identifier)
			*leaderId = uint(id)
		}

		heartBeatServer.OnDisconnect = func(identifier string) {

			err := syscall.Kill(syscall.Getpid(), SignalToTellOfflineNode)
			utility.CheckError(err)

			fmt.Printf("Node %s is OFFLINE\n", identifier)

			server.Close()

		}

		go func() {
			<-timerFirstHeartBeat.C

			fmt.Printf("Node NOT RESPOND!\n")

			err := syscall.Kill(syscall.Getpid(), SignalToTellOfflineNode)
			utility.CheckError(err)

			server.Close()
		}()

		server.ListenAndServe()
	}()
}
