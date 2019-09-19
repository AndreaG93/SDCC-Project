package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/replyregister"
	"math"
)

type AFTReduceTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
}

type AFTReduceTask struct {
	output                      *AFTReduceTaskOutput
	replyRegistry               *replyregister.Register
	replyChannel                chan interface{}
	faultToleranceLevel         int
	requestsSend                int
	workerProcessesRPCAddresses []string
	reduceTaskIndex             int
	reduceTaskIdentifierDigest  string
}

func NewAFTReduceTask(targetNodeGroupId int, reduceTaskIdentifierDigest string, reduceTaskIndex int) *AFTReduceTask {

	output := new(AFTReduceTask)

	(*output).workerProcessesRPCAddresses, _ = process.GetMembershipRegister().GetWorkerProcessPublicInternetAddressesForRPC(targetNodeGroupId, aftmapreduce.WordCountReduceTaskRPCBasePort)

	(*output).replyChannel = make(chan interface{})
	(*output).faultToleranceLevel = int(math.Floor(float64((len((*output).workerProcessesRPCAddresses) - 1) / 2)))
	(*output).replyRegistry = replyregister.New((*output).faultToleranceLevel + 1)

	(*output).output = new(AFTReduceTaskOutput)
	(*(*output).output).IdGroup = targetNodeGroupId

	(*output).requestsSend = 0

	(*output).reduceTaskIndex = reduceTaskIndex
	(*output).reduceTaskIdentifierDigest = reduceTaskIdentifierDigest

	return output
}

func (obj *AFTReduceTask) GetOutput() interface{} {

	digest, nodeIds, _ := (*obj).replyRegistry.GetMostMatchedReply()

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
	return (*obj).replyRegistry.AddReplyCheckingRequiredMatches(reduceLastReply.Digest, reduceLastReply.NodeId, nil)
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
