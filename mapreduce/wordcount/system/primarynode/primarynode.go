package primarynode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/heartbeatleader"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/noderegister"
)

const anyLeaderElected = 0

type PrimaryNode struct {
	id       uint
	leaderId uint
}

func New(id uint) *PrimaryNode {

	output := new(PrimaryNode)

	(*output).id = id
	(*output).leaderId = anyLeaderElected
	return output
}

func (obj *PrimaryNode) isLeader() bool {
	return (*obj).id == (*obj).leaderId
}

func (obj *PrimaryNode) isLeaderNotElected() bool {
	return (*obj).leaderId == anyLeaderElected
}

func (obj *PrimaryNode) StartWork() {

	startNewElection := make(chan bool, 1)
	waitEndElection := make(chan bool, 1)

	go heartbeat.LeaderElection((*obj).id, startNewElection, waitEndElection)

	go heartbeat.Receive((*obj).id, startNewElection, waitEndElection)

	primaryNodeIds := noderegister.GetInstance().GetPrimaryNodeIDs()

	for recipientId := uint(1); recipientId < uint(len(primaryNodeIds)); recipientId++ {
		if recipientId != (*obj).id {
			go heartbeat.Send((*obj).id, recipientId, startNewElection, waitEndElection)
		}
	}

	system.StartAcceptingRPCRequest(&wordcount.Request{}, (*obj).id)
	//utility.CheckError(err)

	/*
		go heartbeatleader.Send((*obj).id, startNewElection, waitEndElection)

		leaderOfflineChannel := make(chan os.Signal, 1)
		signal.Notify(leaderOfflineChannel, network.SignalToTellOfflineNode)

		(*obj).leaderId, err = zookeeper.GetCurrentLeaderId()
		utility.CheckError(err)

		alreadyRegister = false

		(*obj).leaderId, _ = zookeeper.GetCurrentLeaderId()

		for {

			if (*obj).isLeaderNotElected() {

				fmt.Println("My ID ", (*obj).id, "! -> Any leader elected")

				(*obj).leaderId = zookeeper.StartLeaderElection((*obj).id)

				fmt.Println("Actual leader is ", (*obj).leaderId)

			} else if (*obj).isLeader() {

				(*obj).startWorkAsLeader()

			} else {

				fmt.Println("My ID ", (*obj).id, "! -> Leader's ID ", (*obj).leaderId, "!")

				network.ReceiveHeartbeatFromLeaderNode((*obj).id, &(*obj).leaderId, alreadyRegister)
				alreadyRegister = true

				<-leaderOfflineChannel

				fmt.Println("My ID ", (*obj).id, "! -> Leader offline")

				(*obj).leaderId = anyLeaderElected
			}
		}
	*/

}

func (obj *PrimaryNode) startWorkAsLeader() {
	/*
		primaryNodeIds := noderegister.GetInstance().GetPrimaryNodeIDs()

		for recipientId := uint(1); recipientId < uint(len(primaryNodeIds)); recipientId++ {
			if recipientId != (*obj).id {
				network.SendHeartbeatToPrimaryBackup((*obj).id, recipientId)
			}
		}

		system.StartAcceptingRPCRequest(&wordcount.Request{}, (*obj).id) //utility.CheckError(err)
	*/

}
