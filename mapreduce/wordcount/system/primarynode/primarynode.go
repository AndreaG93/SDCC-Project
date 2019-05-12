package primarynode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/heartbeat"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/noderegister"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/timerregister"
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	var err error

	go func() {
		heartbeat.ReceiveHeartBeats((*obj).id, &(*obj).leaderId)
	}()

	leaderOfflineChannel := make(chan os.Signal, 1)
	signal.Notify(leaderOfflineChannel, syscall.SIGUSR1)

	(*obj).leaderId, err = zookeeper.GetCurrentLeaderId()
	utility.CheckError(err)

	for {

		if (*obj).isLeaderNotElected() {

			fmt.Println("My ID ", (*obj).id, "! -> Any leader elected")

			(*obj).leaderId = zookeeper.StartLeaderElection((*obj).id)

			fmt.Println("Actual leader is ", (*obj).leaderId)

		} else if (*obj).isLeader() {

			(*obj).startWorkAsLeader()

		} else {

			fmt.Println("My ID ", (*obj).id, "! -> Leader's ID ", (*obj).leaderId, "!")

			go func() {
				timerregister.GetInstance().StartTimer((*obj).leaderId)

				fmt.Printf("Timer associated to leader %d expired", (*obj).leaderId)

				err = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
				utility.CheckError(err)
			}()

			<-leaderOfflineChannel

			fmt.Println("My ID ", (*obj).id, "! -> Leader offline")

			(*obj).leaderId = anyLeaderElected
		}
	}
}

func (obj *PrimaryNode) startWorkAsLeader() {

	primaryNodeIds := noderegister.GetInstance().GetPrimaryNodeIDs()

	for recipientId := uint(1); recipientId < uint(len(primaryNodeIds)); recipientId++ {
		if recipientId != (*obj).id {
			go heartbeat.SendHeartBeatsTo((*obj).id, recipientId)
		}
	}

	system.StartAcceptingRPCRequest(&wordcount.Request{}, (*obj).id) //utility.CheckError(err)
}
