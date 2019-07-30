package primary

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/implementations/wordcount"
	"SDCC-Project/aftmapreduce/node"
	"encoding/gob"
	"fmt"
)

type Primary struct {
	id                         int
	mapReduceRequestRPCAddress string
	isLeader                   chan bool
}

func New(id int, internetAddress string) *Primary {

	output := new(Primary)

	node.Initialize(id, "Primary")

	(*output).id = id
	(*output).mapReduceRequestRPCAddress = fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRequestRPCBasePort+id)
	(*output).isLeader = make(chan bool)

	return output
}

func (obj *Primary) StartWork() {

	go node.GetZookeeperClient().RunAsLeaderCandidate((*obj).isLeader, (*obj).mapReduceRequestRPCAddress)

	<-(*obj).isLeader

	node.GetLogger().PrintMessage("I'm leader")

	gob.Register(wordcount.Input{})
	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	aftmapreduce.InitNeededZNodePathsToManageClientsRequests()

	for _, item := range aftmapreduce.GetPendingClientsRequests() {
		go aftmapreduce.ManageClientRequest(item)
	}

	aftmapreduce.StartAcceptingRPCRequest(&aftmapreduce.EntryPoint{}, (*obj).mapReduceRequestRPCAddress)
}
