package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/registry/reply"
	"fmt"
	"math"
	"net/rpc"
	"time"
)

type AFTMapTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
	MappedDataSizes          map[int]int
}

type MapTask struct {
	mapTaskOutput       *AFTMapTaskOutput
	workersReplyChannel chan *MapOutput
	faultToleranceLevel int
	requestSend         int
	workersAddresses    []string
	registry            *reply.MapReplyRegistry
	split               string
}

func NewMapTask(split string, workerGroupId int) *MapTask {

	output := new(MapTask)

	(*output).mapTaskOutput = new(AFTMapTaskOutput)
	(*(*output).mapTaskOutput).IdGroup = workerGroupId
	(*output).workersReplyChannel = make(chan *MapOutput)
	(*output).requestSend = 0
	(*output).workersAddresses = node.GetMembershipRegister().GetWorkerProcessPublicInternetAddressesForRPC(workerGroupId, aftmapreduce.WordCountMapTaskRPCBasePort)

	(*output).split = split

	(*output).faultToleranceLevel = int(math.Floor(float64((len((*output).workersAddresses) - 1) / 2)))
	(*output).registry = reply.NewMapReplyRegistry((*output).faultToleranceLevel + 1)

	return output
}

func (obj *MapTask) Execute() *AFTMapTaskOutput {

	defer close((*obj).workersReplyChannel)

	for ; (*obj).requestSend <= (*obj).faultToleranceLevel; (*obj).requestSend++ {
		go executeSingleMapTaskReplica((*obj).split, (*obj).workersAddresses[(*obj).requestSend], (*obj).workersReplyChannel)
	}

	(*obj).startListeningWorkersReplies()

	digest, nodeIds, mappedDataSizes := (*obj).registry.GetMostMatchedReply()

	(*(*obj).mapTaskOutput).ReplayDigest = digest
	(*(*obj).mapTaskOutput).MappedDataSizes = mappedDataSizes
	(*(*obj).mapTaskOutput).NodeIdsWithCorrectResult = nodeIds

	return (*obj).mapTaskOutput
}

func (obj *MapTask) startListeningWorkersReplies() {

	repliesReceived := 0
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			node.GetLogger().PrintInfoTaskMessage("AFT-MAP-TASK", "Timout expired!")

			if (*obj).requestSend < len((*obj).workersAddresses) {
				go executeSingleMapTaskReplica((*obj).split, (*obj).workersAddresses[(*obj).requestSend], (*obj).workersReplyChannel)
				(*obj).requestSend++
			} else {
				panic("number of available WP isn't enough")
			}

		case myReply := <-(*obj).workersReplyChannel:

			timeout.Stop()
			repliesReceived++

			node.GetLogger().PrintInfoTaskMessage("AFT-MAP-TASK", fmt.Sprintf("Received reply by node id %d group %d", myReply.IdNode, myReply.IdGroup))

			if (*obj).registry.Add(myReply.ReplayDigest, myReply.IdNode, myReply.MappedDataSizes) {
				return
			}

			if repliesReceived < (*obj).requestSend {
				timeout.Reset(1 * time.Second)
				continue
			}

			if (*obj).requestSend < len((*obj).workersAddresses) {
				go executeSingleMapTaskReplica((*obj).split, (*obj).workersAddresses[(*obj).requestSend], (*obj).workersReplyChannel)
				(*obj).requestSend++
				timeout.Reset(1 * time.Second)
			} else {
				panic(fmt.Sprintf("number of available WP isn't enough -- Group ID %d", (*obj).mapTaskOutput.IdGroup))
			}
		}
	}
}

func executeSingleMapTaskReplica(split string, fullRPCInternetAddress string, reply chan *MapOutput) {

	input := new(MapInput)
	output := new(MapOutput)

	(*input).Text = split
	(*input).MappingCardinality = node.GetPropertyAsInteger(property.MapCardinality)

	worker, err := rpc.Dial("tcp", fullRPCInternetAddress)
	if err != nil {
		node.GetLogger().PrintErrorTaskMessage(MapTaskName, err.Error())
		return
	}

	err = worker.Call("Map.Execute", input, output)
	if err == nil {
		reply <- output
	}
}
