package workernode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/heartbeat"
	"SDCC-Project-WorkerNode/utility"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
)

type WorkerNode struct {
	id                           uint
	zookeeperServerPoolAddresses []string
}

func New(id uint, zookeeperServerPoolAddresses []string) *WorkerNode {

	output := new(WorkerNode)

	(*output).id = id
	(*output).zookeeperServerPoolAddresses = zookeeperServerPoolAddresses

	return output
}

func (obj *WorkerNode) StartWork() {

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Map{}, (*obj).id); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Reduce{}, (*obj).id); err != nil {
			panic(err)
		}
	}()

	(*obj).startHeartbeat()
}

func (obj *WorkerNode) startHeartbeat() {

	var zkLeaderChangeEventChannel <-chan zk.Event
	var zkConnection *zk.Conn
	var currentLeaderId int
	var data []byte
	var err error

	stopHeartBeat := false

	zkConnection, _, err = zookeeper.ConnectToZookeeperServers((*obj).zookeeperServerPoolAddresses)
	utility.CheckError(err)

	for {

		data, _, zkLeaderChangeEventChannel, err = zkConnection.GetW("/current_leader_id")
		utility.CheckError(err)

		currentLeaderId, err = strconv.Atoi(string(data))
		utility.CheckError(err)

		if currentLeaderId != -1 {
			heartbeat.SendStoppableHeartBeatsTo((*obj).id, uint(currentLeaderId), &stopHeartBeat)
		}

		<-zkLeaderChangeEventChannel
		stopHeartBeat = true
	}
}
