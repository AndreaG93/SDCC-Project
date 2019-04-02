package primarynode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/zookeeper"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"github.com/codeskyblue/heartbeat"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	leaderNotElected = -1
)

type PrimaryNode struct {
	id                      int
	leaderId                int
	listenPortForRPC        string
	listenUDPPort           int
	allPrimaryNodeAddresses []string
}

func New(primaryNodeId int, listenUDPPort int) *PrimaryNode {

	output := new(PrimaryNode)

	(*output).id = primaryNodeId
	(*output).leaderId = leaderNotElected
	(*output).listenUDPPort = listenUDPPort
	(*output).allPrimaryNodeAddresses = []string{"localhost:5000", "localhost:5001", "localhost:5002", "localhost:5003", "localhost:5004"}

	return output
}

func (obj *PrimaryNode) StartWork() {

	go func() {
		if err := system.StartAcceptingRPCRequest(&wordcount.Request{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()
}

func (obj *PrimaryNode) startToSendHeartbeatFromLeader(leaderOfflineEventChannel chan bool) {

	hbs := heartbeat.NewServer("my-secret", 3*time.Second)

	hbs.OnDisconnect = func(identifier string) {
		fmt.Println(identifier, "is offline")
		leaderOfflineEventChannel <- true
	}
	http.Handle("/heartbeatFromLeader", hbs)
	utility.CheckError(http.ListenAndServe(":7000", nil))
}

func (obj *PrimaryNode) startToReceiveHeartbeatFromLeader(leaderOfflineEventChannel chan bool) {

	hbs := heartbeat.NewServer("my-secret", 3*time.Second)

	hbs.OnDisconnect = func(identifier string) {
		fmt.Println(identifier, "is offline")
		leaderOfflineEventChannel <- true
	}
	http.Handle("/heartbeatFromLeader", hbs)
	utility.CheckError(http.ListenAndServe(":7000", nil))
}

func (obj *PrimaryNode) startSendingHeartbeatToBackups() {

	for index := 0; index < len((*obj).allPrimaryNodeAddresses); index++ {

		if index != (*obj).id {
			go func() {
				system.SendHeartbeat((*obj).allPrimaryNodeAddresses[index])
			}()
		}
	}
}

func (obj *PrimaryNode) StartWork() {

	var leaderNotRespondingChannel chan bool

	for {

		if (*obj).leaderId == leaderNotElected {

			fmt.Println("I'm node id ", (*obj).id, ": i disclose any leader!")
			(*obj).leaderId = zookeeper.StartLeaderElection((*obj).id)
			fmt.Println("Actual leader is ", (*obj).leaderId)

		} else if (*obj).leaderId != (*obj).id {

			fmt.Println("I'm node id ", (*obj).id, ": i disclose leader id ", (*obj).leaderId)

			go system.ReceivingHeartbeat(leaderNotRespondingChannel, (*obj).listenUDPPort)

			<-leaderNotRespondingChannel

			fmt.Println("I'm node id ", (*obj).id, ": leader doesn't respond")

			(*obj).leaderId = leaderNotElected

		} else {

			(*obj).startSendingHeartbeatToBackups()

			/*
				go func() {
					err := system.StartAcceptingRPCRequest(&wordcount.Request{}, (*obj).listenPortForRPC)
					utility.CheckError(err)
				}()

			*/
		}
	}
}
