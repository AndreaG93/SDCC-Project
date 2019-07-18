package Worker

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"fmt"
	"sync"
)

type Worker struct {
	id int
}

func New(id int) *Worker {

	output := new(Worker)
	(*output).id = id

	return output
}

func (obj *Worker) StartWork() {

	go system.StartAcceptingRPCRequest(wordcount.Map{}, (*obj).id)
	go system.StartAcceptingRPCRequest(wordcount.Map{}, (*obj).id)

	(*obj).startToSendHeartbeatToLeader()
}

func (obj *Worker) startToSendHeartbeatToLeader() {

	var myWaitGroup sync.WaitGroup

	myWaitGroup.Add(1)

	go func() {
		fmt.Print("Dsadasd")
	}()

	myWaitGroup.Wait()

	/*
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

	*/
}
