package Task

import (
	"SDCC-Project/MapReduce/Input"
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
	faultToleranceLevel int
	currentWorkerID     int
	workersAddresses    []string
	repliesRegistry     *WorkersResponsesRegistry.WorkersResponsesRegistry
	inputSplit          Input.MiddleInput
}

func NewBFTMapReduce(inputSplit Input.MiddleInput, faultToleranceLevel int, workersAddresses []string) *BFTMapOrReduceService {

	output := new(BFTMapOrReduceService)

	(*output).workersReplyChannel = make(chan workerReply)
	(*output).faultToleranceLevel = faultToleranceLevel
	(*output).currentWorkerID = 0
	(*output).workersAddresses = workersAddresses
	(*output).repliesRegistry = WorkersResponsesRegistry.New((*output).faultToleranceLevel + 1)
	(*output).inputSplit = inputSplit

	return output
}

/**
 * This function is used to execute a 'Byzantine Fault Tolerant' Map-Task.
 */
func (obj *BFTMapOrReduceService) Execute() (string, []string) {

	defer close((*obj).workersReplyChannel)

	numberOfReply := 0

	for ; (*obj).currentWorkerID <= (*obj).faultToleranceLevel; (*obj).currentWorkerID++ {
		go executeSingleMapTaskReplica((*obj).inputSplit, (*obj).workersAddresses[(*obj).currentWorkerID], (*obj).workersReplyChannel)
	}

	for reply := range (*obj).workersReplyChannel {

		if (*obj).repliesRegistry.AddWorkerResponse(reply.digest, reply.workerAddress) {
			break
		}

		numberOfReply++

		if numberOfReply >= (*obj).faultToleranceLevel+1 {
			(*obj).currentWorkerID++
			go executeSingleMapTaskReplica((*obj).inputSplit, (*obj).workersAddresses[(*obj).currentWorkerID], (*obj).workersReplyChannel)
		}
	}

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
			(*obj).currentWorkerID++
			go executeSingleMapTaskReplica((*obj).inputSplit, (*obj).workersAddresses[(*obj).currentWorkerID], (*obj).workersReplyChannel)
		}
	}
}

func executeSingleMapTaskReplica(middleInput Input.MiddleInput, address string, reply chan workerReply) {

	var workerInput MapReduceInput
	var workerOutput MapReduceOutput

	worker, err := rpc.Dial("tcp", address)
	utility.CheckError(err)

	workerInput.InputData = middleInput

	err = worker.Call("MapReduce.Execute", &workerInput, &workerOutput)

	if err == nil {

		output := new(workerReply)
		(*output).workerAddress = address
		(*output).digest = workerOutput.Digest

		reply <- *output
	}
}
