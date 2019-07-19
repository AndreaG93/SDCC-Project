package BFTMapTask

import (
	"SDCC-Project/MapReduce/Data"
	"SDCC-Project/MapReduce/Registry/WorkersResponsesRegistry"
	"SDCC-Project/MapReduce/wordcount"
	"SDCC-Project/utility"
	"net/rpc"
)

type workerReply struct {
	digest        string
	workerAddress string
}

type BFTMapTask struct {
	workersReplyChannel chan workerReply
	killTaskChannel     chan bool
	faultToleranceLevel int
	currentWorkerID     int
	workersAddresses    []string
	repliesRegistry     *WorkersResponsesRegistry.WorkersResponsesRegistry
	inputSplit          Data.Split
}

func New(inputSplit Data.Split, faultToleranceLevel int, workersAddresses []string) *BFTMapTask {

	output := new(BFTMapTask)

	(*output).workersReplyChannel = make(chan workerReply)
	(*output).killTaskChannel = make(chan bool)
	(*output).faultToleranceLevel = faultToleranceLevel
	(*output).currentWorkerID = 0
	(*output).workersAddresses = workersAddresses
	(*output).repliesRegistry = WorkersResponsesRegistry.New((*output).faultToleranceLevel, (*output).killTaskChannel)
	(*output).inputSplit = inputSplit

	return output
}

/**
 * This function is used to execute a 'Byzantine Fault Tolerant' Map-Task.
 */
func (obj *BFTMapTask) Execute() (string, []string) {

	defer close((*obj).workersReplyChannel)
	defer close((*obj).killTaskChannel)

	go (*obj).startListeningWorkersReplies()

	for ; (*obj).currentWorkerID <= (*obj).faultToleranceLevel; (*obj).currentWorkerID++ {
		go (*obj).executeSingleMapTaskReplica()
	}

	<-(*obj).killTaskChannel
	return (*obj).repliesRegistry.GetMostMatchedWorkerResponse()
}

func (obj *BFTMapTask) startListeningWorkersReplies() {

	numberOfReply := 0

	for response := range (*obj).workersReplyChannel {

		if (*obj).repliesRegistry.AddWorkerResponse(response.digest, response.workerAddress) {
			return
		}

		numberOfReply++

		if numberOfReply >= (*obj).faultToleranceLevel+1 {
			go (*obj).executeSingleMapTaskReplica()
		}
	}
}

func (obj *BFTMapTask) executeSingleMapTaskReplica() {

	var mapTaskInput wordcount.MapInput
	var mapTaskOutput wordcount.MapOutput

	address := (*obj).workersAddresses[(*obj).currentWorkerID]

	worker, err := rpc.Dial("tcp", address)
	utility.CheckError(err)

	mapTaskInput.Input = (*obj).inputSplit
	mapTaskInput.MapCardinality = 5

	err = worker.Call("Map.Execute", &mapTaskInput, &mapTaskOutput)

	if err != nil {

		output := new(workerReply)
		(*output).workerAddress = address
		(*output).digest = mapTaskOutput.Digest

		(*obj).workersReplyChannel <- *output
	}
}
