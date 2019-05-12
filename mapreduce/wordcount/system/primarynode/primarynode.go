package primarynode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/heartbeat"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/nodeidsregister"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	leaderNotElected = -1
)

type PrimaryNode struct {
	id                      uint
	leaderId                uint
	listenPortForRPC        string
	listenUDPPort           int
}

func New(primaryNodeId uint, ì) *PrimaryNode {

	output := new(PrimaryNode)

	(*output).id = primaryNodeId
	(*output).leaderId = leaderNotElected

	return output
}

func (obj *PrimaryNode) StartWork(onlinePrimaryNodeRegister map[uint]bool, onlineWorkerNodeRegister map[uint]bool) {

	go func() {
		if err := system.StartAcceptingRPCRequest(&wordcount.Request{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()

	go func() {
		heartbeat.WorkerMonitoring()
	}()

	go func() {
		(*obj).leaderManagement()
	}()
}

func (obj *PrimaryNode) leaderManagement() {

	leaderNotRespondingChannel := make(chan os.Signal, 1)
	signal.Notify(leaderNotRespondingChannel, syscall.SIGUSR1)

	for {

		if (*obj).leaderId == leaderNotElected {

			fmt.Println("I'm node id ", (*obj).id, ": i disclose any leader!")
			(*obj).leaderId = zookeeper.StartLeaderElection((*obj).id)
			fmt.Println("Actual leader is ", (*obj).leaderId)

		} else if (*obj).leaderId != (*obj).id {

			fmt.Println("I'm node id ", (*obj).id, ": i disclose leader id ", (*obj).leaderId)

			go heartbeat.LeaderMonitoring((*obj).leaderId)

			<-leaderNotRespondingChannel

			fmt.Println("I'm node id ", (*obj).id, ": leader doesn't respond")

			(*obj).leaderId = leaderNotElected

		} else {

			for id := range nodeidsregister.GetInstance().GetPrimaryNodeIDs() {

				if uint(id) != (*obj).id {
					heartbeat.StartHeartBeating((*obj).id, nodeidsregister.GetInstance().GetNodeIpAddress(uint(id)))
				}

				done := make(chan bool, 1)
				<-done
			}
		}
	}
}
