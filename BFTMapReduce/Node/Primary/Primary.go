package Primary

import (
	"SDCC-Project/BFTMapReduce"
	"SDCC-Project/BFTMapReduce/ConcreteImplementations/WordCount"
	"SDCC-Project/BFTMapReduce/Task"
	"SDCC-Project/BFTMapReduce/clientrequest"
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
	(*output).mapReduceRequestRPCAddress = fmt.Sprintf("%s:%d", internetAddress, BFTMapReduce.MapReduceRequestRPCBasePort+id)
	(*output).zookeeperClient = zookeeper.New([]string{"localhost:2181"})
	(*output).isLeader = make(chan bool)

	return output
}

func (obj *Primary) StartWork() {

	go (*obj).zookeeperClient.RunAsLeaderCandidate((*obj).isLeader)

	<-(*obj).isLeader

	fmt.Println("I'm LEADER")

	gob.Register(WordCount.File{})
	gob.Register(WordCount.MapInput{})
	gob.Register(WordCount.ReduceInput{})

	clientrequest.InitializationClientsRequestsPath((*obj).zookeeperClient)

	BFTMapReduce.StartAcceptingRPCRequest(&Task.MapReduceRequest{}, (*obj).mapReduceRequestRPCAddress)
}
