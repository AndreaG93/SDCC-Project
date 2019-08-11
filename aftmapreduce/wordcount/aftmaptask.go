package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/registry"
	"SDCC-Project/utility"
	"net/rpc"
)

type MapTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
	MappedDataSizes          map[int]int
}

type MapTask struct {
	mapTaskOutput       *MapTaskOutput
	workersReplyChannel chan *MapOutput
	faultToleranceLevel int
	requestSend         int
	workersAddresses    []string
	registry            *registry.MapReplies
	split               string
}

func NewMapTask(split string, workerGroupId int) *MapTask {

	output := new(MapTask)

	(*output).mapTaskOutput = new(MapTaskOutput)
	(*(*output).mapTaskOutput).IdGroup = workerGroupId
	(*output).workersReplyChannel = make(chan *MapOutput)
	(*output).faultToleranceLevel = 1
	(*output).requestSend = 0
	(*output).workersAddresses = node.GetZookeeperClient().GetWorkerInternetAddressesForRPC(workerGroupId, aftmapreduce.WordCountMapTaskRPCBasePort)
	(*output).registry = registry.NewMapReply((*output).faultToleranceLevel + 1)
	(*output).split = split

	return output
}

func (obj *MapTask) Execute() *MapTaskOutput {

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

	numberOfReply := 0

	for reply := range (*obj).workersReplyChannel {

		if (*obj).registry.Add(reply.ReplayDigest, reply.IdNode, reply.MappedDataSizes) {
			return
		}

		numberOfReply++

		if numberOfReply >= (*obj).faultToleranceLevel+1 {
			(*obj).requestSend++
			go executeSingleMapTaskReplica((*obj).split, (*obj).workersAddresses[(*obj).requestSend], (*obj).workersReplyChannel)
		}
	}
}

func executeSingleMapTaskReplica(split string, fullRPCInternetAddress string, reply chan *MapOutput) {

	input := new(MapInput)
	output := new(MapOutput)

	(*input).Text = split

	worker, err := rpc.Dial("tcp", fullRPCInternetAddress)
	utility.CheckError(err)

	err = worker.Call("Map.Execute", input, output)
	if err == nil {
		reply <- output
	}
}
