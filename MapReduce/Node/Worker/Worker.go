package Worker

import (
	"SDCC-Project/MapReduce"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/MapReduce/Task"
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/utility"
	"encoding/gob"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
)

type Worker struct {
	id                     int
	mapReduceRPCAddress    string
	mapReduceGetRPCAddress string
}

func New(id int, mapReduceRPCAddress string, mapReduceGetRPCAddress string) *Worker {

	output := new(Worker)
	(*output).id = id
	(*output).mapReduceRPCAddress = mapReduceRPCAddress
	(*output).mapReduceGetRPCAddress = mapReduceGetRPCAddress

	return output
}

func (obj *Worker) StartWork() {

	gob.Register(WordCount.MapInput{})
	gob.Register(WordCount.ReduceInput{})

	go MapReduce.StartAcceptingRPCRequest(&Task.MapReduce{}, (*obj).mapReduceRPCAddress)
	go MapReduce.StartAcceptingRPCRequest(&Task.MapReduceGet{}, (*obj).mapReduceGetRPCAddress)

	(*obj).startToSendHeartbeatToLeader()
}

func (obj *Worker) startToSendHeartbeatToLeader() {

	stopHeartBeat := false
	zookeeperClient := zookeeper.New([]string{"localhost:2181"})

	for {

		data, zkLeaderChangeEventChannel := zookeeperClient.GetZNodeData(zookeeper.ActualLeaderZNodePath)

		currentLeaderId, err := strconv.Atoi(string(data))
		utility.CheckError(err)

		if currentLeaderId != -1 {
			heartbeat.SendStoppableHeartBeatsTo((*obj).id, uint(currentLeaderId), &stopHeartBeat)
		}

		<-zkLeaderChangeEventChannel
		stopHeartBeat = true
	}
}
