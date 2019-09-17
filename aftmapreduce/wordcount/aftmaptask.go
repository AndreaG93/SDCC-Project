package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/registry/reply"
	"math"
)

type AFTMapTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
	MappedDataSizes          map[int]int
}

type MapTask struct {
	mapTaskOutput                       *AFTMapTaskOutput
	workersReplyChannel                 chan interface{}
	firstReplyPredictedAsCorrectChannel chan interface{}
	faultToleranceLevel                 int
	requestSend                         int
	workersAddresses                    []string
	registry                            *reply.MapReplyRegistry
	split                               string
}

func NewMapTask(split string, workerGroupId int, firstReplyPredictedAsCorrectChannel chan interface{}) *MapTask {

	output := new(MapTask)

	(*output).mapTaskOutput = new(AFTMapTaskOutput)
	(*(*output).mapTaskOutput).IdGroup = workerGroupId
	(*output).workersReplyChannel = make(chan interface{})
	(*output).workersAddresses, _ = node.GetMembershipRegister().GetWorkerProcessPublicInternetAddressesForRPC(workerGroupId, aftmapreduce.WordCountMapTaskRPCBasePort)
	(*output).split = split
	(*output).faultToleranceLevel = int(math.Floor(float64((len((*output).workersAddresses) - 1) / 2)))
	(*output).registry = reply.NewMapReplyRegistry((*output).faultToleranceLevel + 1)
	(*output).firstReplyPredictedAsCorrectChannel = firstReplyPredictedAsCorrectChannel

	return output
}

func (obj *MapTask) GetOutput() interface{} {

	digest, nodeIds, mappedDataSizes := (*obj).registry.GetMostMatchedReply()

	(*(*obj).mapTaskOutput).ReplayDigest = digest
	(*(*obj).mapTaskOutput).MappedDataSizes = mappedDataSizes
	(*(*obj).mapTaskOutput).NodeIdsWithCorrectResult = nodeIds

	return (*obj).mapTaskOutput
}

func (obj *MapTask) GetReplyChannel() chan interface{} {
	return (*obj).workersReplyChannel
}

func (obj *MapTask) GetFaultToleranceLevel() int {
	return (*obj).faultToleranceLevel
}

func (obj *MapTask) GetAvailableWorkerProcessesRPCInternetAddresses() []string {
	return (*obj).workersAddresses
}

func (obj *MapTask) DoWeHaveEnoughMatchingReplyAfter(lastReply interface{}) bool {
	mapLastReply := lastReply.(*MapOutput)
	return (*obj).registry.Add(mapLastReply.ReplayDigest, mapLastReply.IdNode, mapLastReply.MappedDataSizes)
}

func (obj *MapTask) ExecuteRPCCallTo(fullRPCInternetAddress string) {

	replyChannel := (*obj).workersReplyChannel

	input := MapInput{
		Text:               (*obj).split,
		MappingCardinality: (*node.GetMembershipRegister()).GetGroupAmount(),
	}
	output := MapOutput{}

	go aftmapreduce.MakeRPCCall("Map.Execute", fullRPCInternetAddress, input, &output, replyChannel)
}

func (obj *MapTask) GetChannelToSendFirstReplyPredictedAsCorrect() chan interface{} {
	return (*obj).firstReplyPredictedAsCorrectChannel
}
