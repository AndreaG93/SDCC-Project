package Task

import (
	"SDCC-Project/MapReduce/Data"
	"SDCC-Project/MapReduce/Registry/WorkersResponsesRegistry"
	"SDCC-Project/utility"
	"net/rpc"
)

type workerReply struct {
	digest        string
	workerAddress string
}

type BFTMapOrReduceService struct {
	workersReplyChannel chan workerReply
	killTaskChannel     chan bool
	faultToleranceLevel int
	currentWorkerID     int
	workersAddresses    []string
	repliesRegistry     *WorkersResponsesRegistry.WorkersResponsesRegistry
	inputSplit          Data.Split
}

func NewBFTMapReduce(inputSplit Data.Split, faultToleranceLevel int, workersAddresses []string) *BFTMapOrReduceService {

	output := new(BFTMapOrReduceService)

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
func (obj *BFTMapOrReduceService) Execute() (string, []string) {

	defer close((*obj).workersReplyChannel)
	defer close((*obj).killTaskChannel)

	go (*obj).startListeningWorkersReplies()

	for ; (*obj).currentWorkerID <= (*obj).faultToleranceLevel; (*obj).currentWorkerID++ {
		go (*obj).executeSingleMapTaskReplica()
	}

	<-(*obj).killTaskChannel
	return (*obj).repliesRegistry.GetMostMatchedWorkerResponse()
}

func (obj *BFTMapOrReduceService) startListeningWorkersReplies() {

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

func (obj *BFTMapOrReduceService) executeSingleMapTaskReplica() {

	var input MapReduceInput
	var output MapReduceOutput

	address := (*obj).workersAddresses[(*obj).currentWorkerID]

	worker, err := rpc.Dial("tcp", address)
	utility.CheckError(err)

	input.InputData = (*obj).inputSplit

	err = worker.Call("Map.Execute", &input, &output)

	if err != nil {

		output := new(workerReply)
		(*output).workerAddress = address
		(*output).digest = output.digest

		(*obj).workersReplyChannel <- *output
	}
}
