package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/replyregister"
	"fmt"
	"math"
)

type AFTMapTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
	MappedDataSizes          map[int]int
}

type AFTMapTask struct {
	mapTaskOutput                       *AFTMapTaskOutput
	workersReplyChannel                 chan interface{}
	firstReplyPredictedAsCorrectChannel chan interface{}
	faultToleranceLevel                 int
	requestSend                         int
	workersAddresses                    []string
	registry                            *replyregister.Register
	split                               string
}

func NewAFTMapTask(split string, workerGroupId int, firstReplyPredictedAsCorrectChannel chan interface{}) *AFTMapTask {

	output := new(AFTMapTask)

	(*output).mapTaskOutput = new(AFTMapTaskOutput)
	(*(*output).mapTaskOutput).IdGroup = workerGroupId
	(*output).workersReplyChannel = make(chan interface{})
	(*output).workersAddresses, _ = process.GetMembershipRegister().GetWorkerProcessPublicInternetAddressesForRPC(workerGroupId, aftmapreduce.WordCountMapTaskRPCBasePort)
	(*output).split = split
	(*output).faultToleranceLevel = int(math.Floor(float64((len((*output).workersAddresses) - 1) / 2)))
	(*output).registry = replyregister.New((*output).faultToleranceLevel + 1)
	(*output).firstReplyPredictedAsCorrectChannel = firstReplyPredictedAsCorrectChannel

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("An AFT Map-Task is STARTING with FT %d :: WPG %d :: available WPs: %s", (*output).faultToleranceLevel, workerGroupId, (*output).workersAddresses))

	return output
}

func (obj *AFTMapTask) GetOutput() interface{} {

	digest, nodeIds, mappedDataSizes := (*obj).registry.GetMostMatchedReply()

	(*(*obj).mapTaskOutput).ReplayDigest = digest
	(*(*obj).mapTaskOutput).MappedDataSizes = mappedDataSizes.(map[int]int)
	(*(*obj).mapTaskOutput).NodeIdsWithCorrectResult = nodeIds

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("An AFT Map-Task is COMPLETE with digest %s :: WPs with correct results: %d :: IndexPartitionSize %d", digest, nodeIds, mappedDataSizes.(map[int]int)))

	return (*obj).mapTaskOutput
}

func (obj *AFTMapTask) GetReplyChannel() chan interface{} {
	return (*obj).workersReplyChannel
}

func (obj *AFTMapTask) GetFaultToleranceLevel() int {
	return (*obj).faultToleranceLevel
}

func (obj *AFTMapTask) GetAvailableWorkerProcessesRPCInternetAddresses() []string {
	return (*obj).workersAddresses
}

func (obj *AFTMapTask) DoWeHaveEnoughMatchingReplyAfter(lastReply interface{}) bool {
	mapLastReply := lastReply.(*MapOutput)

	process.GetMembershipRegister().AddProcessCPUUtilization(mapLastReply.MyInternetAddress, mapLastReply.CPUUtilization)
	return (*obj).registry.AddReplyCheckingRequiredMatches(mapLastReply.ReplayDigest, mapLastReply.IdNode, mapLastReply.MappedDataSizes)
}

func (obj *AFTMapTask) ExecuteRPCCallTo(fullRPCInternetAddress string) {

	replyChannel := (*obj).workersReplyChannel

	input := MapInput{
		Text:               (*obj).split,
		MappingCardinality: (*process.GetMembershipRegister()).GetGroupAmount(),
	}
	output := MapOutput{}

	go func() {
		if err := aftmapreduce.MakeRPCCall("Map.Execute", fullRPCInternetAddress, input, &output, replyChannel); err != nil {
			process.GetLogger().PrintInfoLevelMessage(err.Error())
		}
	}()
}

func (obj *AFTMapTask) GetChannelToSendFirstReply() chan interface{} {
	return (*obj).firstReplyPredictedAsCorrectChannel
}
