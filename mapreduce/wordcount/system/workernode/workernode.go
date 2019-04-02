package workernode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/utility"
	"context"
	"github.com/codeskyblue/heartbeat"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

type WorkerNode struct {
	id                            uint
	listenPortForRPC              string
	zookeeperServerPoolAddresses  []string
	heartbeatClient               *heartbeat.Client
	heartbeatClientCancelFunction context.CancelFunc
}

func New(id uint, listenPortForRPC string, zookeeperServerPoolAddresses []string) *WorkerNode {

	output := new(WorkerNode)

	(*output).id = id
	(*output).zookeeperServerPoolAddresses = zookeeperServerPoolAddresses
	(*output).listenPortForRPC = listenPortForRPC
	(*output).heartbeatClient = &heartbeat.Client{
		ServerAddr: "",
		Secret:     "my-secret",
		Identifier: string((*output).id),
	}

	return output
}

func (obj *WorkerNode) StartWork() {

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Map{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Reduce{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()

	(*obj).startHeartbeatService()
}

func (obj *WorkerNode) startHeartbeat() {

	var err error
	var leaderAddress string

	leaderAddress, err = zookeeper.GetCurrentLeaderIPAddress()
	utility.CheckError(err)

	(*obj).heartbeatClient.ServerAddr = leaderAddress
	(*obj).heartbeatClientCancelFunction = (*obj).heartbeatClient.Beat(500 * time.Millisecond)
}

func (obj *WorkerNode) stopHeartbeat() {

	if (*obj).heartbeatClientCancelFunction != nil {
		(*obj).heartbeatClientCancelFunction()
	}
}

func (obj *WorkerNode) startHeartbeatService() {

	var zkLeaderChangeEventChannel <-chan zk.Event
	var zkConnection *zk.Conn
	var currentLeaderId int
	var data []byte
	var err error

	zkConnection, _, err = zookeeper.ConnectToZookeeperServers((*obj).zookeeperServerPoolAddresses)
	utility.CheckError(err)

	for {

		data, _, zkLeaderChangeEventChannel, err = zkConnection.GetW("/current_leader_id")
		utility.CheckError(err)

		currentLeaderId, err = strconv.Atoi(string(data))
		utility.CheckError(err)

		if currentLeaderId != -1 {
			(*obj).startHeartbeat()
		}

		<-zkLeaderChangeEventChannel
		(*obj).stopHeartbeat()
	}
}
