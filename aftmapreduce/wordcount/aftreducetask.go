package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/registry/reply"
	"math"
)

type AFTReduceTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
}

type AFTReduceTask struct {
	output                      *AFTReduceTaskOutput
	replyRegistry               *reply.ReduceReplyRegistry
	replyChannel                chan interface{}
	faultToleranceLevel         int
	requestsSend                int
	workerProcessesRPCAddresses []string
	reduceTaskIndex             int
	reduceTaskIdentifierDigest  string
}

func NewAFTReduceTask(targetNodeGroupId int, reduceTaskIdentifierDigest string, reduceTaskIndex int) *AFTReduceTask {

	output := new(AFTReduceTask)

	(*output).workerProcessesRPCAddresses, _ = node.GetMembershipRegister().GetWorkerProcessPublicInternetAddressesForRPC(targetNodeGroupId, aftmapreduce.WordCountReduceTaskRPCBasePort)

	(*output).replyChannel = make(chan interface{})
	(*output).faultToleranceLevel = int(math.Floor(float64((len((*output).workerProcessesRPCAddresses) - 1) / 2)))
	(*output).replyRegistry = reply.NewReduceReplyRegistry((*output).faultToleranceLevel + 1)

	(*output).output = new(AFTReduceTaskOutput)
	(*(*output).output).IdGroup = targetNodeGroupId

	(*output).requestsSend = 0

	(*output).reduceTaskIndex = reduceTaskIndex
	(*output).reduceTaskIdentifierDigest = reduceTaskIdentifierDigest

	return output
}

func (obj *AFTReduceTask) GetOutput() interface{} {

	digest, nodeIds := (*obj).replyRegistry.GetMostMatchedReply()

	(*(*obj).output).ReplayDigest = digest
	(*(*obj).output).NodeIdsWithCorrectResult = nodeIds

	return (*obj).output
}

func (obj *AFTReduceTask) GetReplyChannel() chan interface{} {
	return (*obj).replyChannel
}

func (obj *AFTReduceTask) GetFaultToleranceLevel() int {
	return (*obj).faultToleranceLevel
}

func (obj *AFTReduceTask) GetAvailableWorkerProcessesRPCInternetAddresses() []string {
	return (*obj).workerProcessesRPCAddresses
}

func (obj *AFTReduceTask) DoWeHaveEnoughMatchingReplyAfter(lastReply interface{}) bool {
	reduceLastReply := lastReply.(*ReduceOutput)
	return (*obj).replyRegistry.Add(reduceLastReply.Digest, reduceLastReply.NodeId)
}

func (obj *AFTReduceTask) ExecuteRPCCallTo(fullRPCInternetAddress string) {

	replyChannel := (*obj).replyChannel

	input := ReduceInput{
		LocalDataDigest: (*obj).reduceTaskIdentifierDigest,
		ReduceWorkIndex: (*obj).reduceTaskIndex,
	}
	output := ReduceOutput{}

	go aftmapreduce.MakeRPCCall("Reduce.Execute", fullRPCInternetAddress, input, &output, replyChannel)
}

func (obj *AFTReduceTask) GetChannelToSendFirstReplyPredictedAsCorrect() chan interface{} {
	return nil
}
