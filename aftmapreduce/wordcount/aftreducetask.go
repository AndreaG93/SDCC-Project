package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/registry/reply"
	"fmt"
	"math"
	"net/rpc"
	"time"
)

type AFTReduceTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
}

type AFTReduceTask struct {
	output *AFTReduceTaskOutput

	replyRegistry *reply.ReduceReplyRegistry
	replyChannel  chan *ReduceOutput

	arbitraryFaultToleranceLevel        int
	requestsSend                        int
	targetNodesFullRPCInternetAddresses []string

	reduceTaskIndex            int
	reduceTaskIdentifierDigest string
}

func NewAFTReduceTask(targetNodeGroupId int, reduceTaskIdentifierDigest string, reduceTaskIndex int) *AFTReduceTask {

	output := new(AFTReduceTask)

	(*output).targetNodesFullRPCInternetAddresses = node.GetMembershipRegister().GetWorkerProcessPublicInternetAddressesForRPC(targetNodeGroupId, aftmapreduce.WordCountReduceTaskRPCBasePort)

	(*output).replyChannel = make(chan *ReduceOutput)
	(*output).arbitraryFaultToleranceLevel = int(math.Floor(float64((len((*output).targetNodesFullRPCInternetAddresses) - 1) / 2)))
	(*output).replyRegistry = reply.NewReduceReplyRegistry((*output).arbitraryFaultToleranceLevel + 1)

	(*output).output = new(AFTReduceTaskOutput)
	(*(*output).output).IdGroup = targetNodeGroupId

	(*output).requestsSend = 0

	(*output).reduceTaskIndex = reduceTaskIndex
	(*output).reduceTaskIdentifierDigest = reduceTaskIdentifierDigest

	return output
}

func (obj *AFTReduceTask) Execute() *AFTReduceTaskOutput {

	defer close((*obj).replyChannel)

	for ; (*obj).requestsSend <= (*obj).arbitraryFaultToleranceLevel; (*obj).requestsSend++ {
		go executeSingleReduceTaskReplica((*obj).reduceTaskIdentifierDigest, (*obj).reduceTaskIndex, (*obj).targetNodesFullRPCInternetAddresses[(*obj).requestsSend], (*obj).replyChannel)
	}

	(*obj).startListeningWorkersReplies()

	digest, nodeIds := (*obj).replyRegistry.GetMostMatchedReply()

	(*(*obj).output).ReplayDigest = digest
	(*(*obj).output).NodeIdsWithCorrectResult = nodeIds

	return (*obj).output
}

func (obj *AFTReduceTask) startListeningWorkersReplies() {

	repliesReceived := 0
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			node.GetLogger().PrintInfoTaskMessage("AFT-MAP-TASK", "Timout expired!")

			if (*obj).requestsSend < len((*obj).targetNodesFullRPCInternetAddresses) {
				go executeSingleReduceTaskReplica((*obj).reduceTaskIdentifierDigest, (*obj).reduceTaskIndex, (*obj).targetNodesFullRPCInternetAddresses[(*obj).requestsSend], (*obj).replyChannel)
				(*obj).requestsSend++
			} else {
				panic(fmt.Sprintf("number of available WP isn't enough -- Group ID %d", (*obj).output.IdGroup))
			}

		case myReply := <-(*obj).replyChannel:

			timeout.Stop()

			if (*obj).replyRegistry.Add(myReply.Digest, myReply.NodeId) {
				return
			}

			if repliesReceived < (*obj).requestsSend {
				timeout.Reset(1 * time.Second)
				continue
			}

			if (*obj).requestsSend < len((*obj).targetNodesFullRPCInternetAddresses) {
				go executeSingleReduceTaskReplica((*obj).reduceTaskIdentifierDigest, (*obj).reduceTaskIndex, (*obj).targetNodesFullRPCInternetAddresses[(*obj).requestsSend], (*obj).replyChannel)
				(*obj).requestsSend++
				timeout.Reset(1 * time.Second)
			} else {
				panic(fmt.Sprintf("number of available WP isn't enough -- Group ID %d", (*obj).output.IdGroup))
			}
		}
	}
}

func executeSingleReduceTaskReplica(localDataDigest string, ReduceWorkIndex int, fullRPCInternetAddress string, reply chan *ReduceOutput) {

	node.GetLogger().PrintInfoTaskMessage("SINGLE-REDUCE-REPLICA", fmt.Sprintf("Send a 'REDUCE' command to: %s -- Reduce Task Index %d", fullRPCInternetAddress, ReduceWorkIndex))

	input := new(ReduceInput)
	output := new(ReduceOutput)

	(*input).LocalDataDigest = localDataDigest
	(*input).ReduceWorkIndex = ReduceWorkIndex

	worker, err := rpc.Dial("tcp", fullRPCInternetAddress)
	if err != nil {
		node.GetLogger().PrintErrorTaskMessage(ReduceTaskName, err.Error())
		return
	}

	err = worker.Call("Reduce.Execute", input, output)
	if err == nil {
		reply <- output
	}
}
