package primary

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/ConcreteImplementations/wordcount"
	"SDCC-Project/aftmapreduce/registries/zookeeperclient"
	"SDCC-Project/cloud/zookeeper"
	"encoding/gob"
	"fmt"
)

type Primary struct {
	id                         int
	mapReduceRequestRPCAddress string
	zookeeperClient            *zookeeper.Client
	isLeader                   chan bool
}

func New(id int, internetAddress string) *Primary {

	output := new(Primary)
	(*output).id = id
	(*output).mapReduceRequestRPCAddress = fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRequestRPCBasePort+id)
	(*output).zookeeperClient = zookeeperclient.GetInstance()
	(*output).isLeader = make(chan bool)

	return output
}

func (obj *Primary) StartWork() {

	go (*obj).zookeeperClient.RunAsLeaderCandidate((*obj).isLeader)

	<-(*obj).isLeader

	fmt.Println("I'm LEADER")

	gob.Register(wordcount.Input{})
	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	aftmapreduce.InitNeededZNodePathsToManageClientsRequests((*obj).zookeeperClient)
	pendingClientsRequests := aftmapreduce.GetPendingClientsRequests((*obj).zookeeperClient)

	for _, item := range pendingClientsRequests {
		go aftmapreduce.ManageClientRequest(item)
	}

	aftmapreduce.StartAcceptingRPCRequest(&aftmapreduce.EntryPoint{}, (*obj).mapReduceRequestRPCAddress)
}
