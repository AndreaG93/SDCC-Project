package heartbeat

import (
	"fmt"
	"github.com/codeskyblue/heartbeat"
	"github.com/samuel/go-zookeeper/zk"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	heartBeatPort         = "30000"
	heartBeatingFrequency = 600 * time.Millisecond
	heartBeatingTimeOut   = 2 * time.Second
	httpEntry             = "/leader"
)

func ReceiveHeartBeatFromWorkers() {

	heartBeatServer := heartbeat.NewServer("my-secret", heartBeatingTimeOut)

	heartBeatServer.OnConnect = func(identifier string, r *http.Request) {
		fmt.Println("Node", identifier, "is ONLINE")
	}

	heartBeatServer.OnDisconnect = func(identifier string) {
		fmt.Println("Node", identifier, "is OFFLINE")
	}

	fmt.Println("Worker heartbeat listening...")

	http.Handle(httpEntry, heartBeatServer)
	go http.ListenAndServe(":"+heartBeatPort, nil)
}

func SendHeartBeatToLeader(senderNodeId uint) {

	var leaderAddress string
	var leaderWatcher <-chan zk.Event

	for {

		//leaderAddress, leaderWatcher = zookeeper.GetActualClusterLeaderAddress()

		result := strings.Split(leaderAddress, ":")

		client := &heartbeat.Client{
			ServerAddr: fmt.Sprintf("http://%s:%s%s", result[0], heartBeatPort, httpEntry),
			Secret:     "my-secret",
			Identifier: strconv.Itoa(int(senderNodeId)),
		}

		cancel := client.Beat(heartBeatingFrequency)

		<-leaderWatcher

		fmt.Println("Leader has been changed!")

		cancel()
	}
}
