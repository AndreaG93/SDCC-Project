package aftmapreduce

import (
	"SDCC-Project/aftmapreduce/registries/WorkersResponsesRegistry"
	"SDCC-Project/utility"
	"net/rpc"
)

type workerReply struct {
	digest        string
	workerAddress string
}

type Task struct {
	workersReplyChannel chan workerReply
	faultToleranceLevel int
	currentWorkerID     int
	workersAddresses    []string
	repliesRegistry     *WorkersResponsesRegistry.WorkersResponsesRegistry
	inputSplit          TransientData
}

func NewTask(inputSplit TransientData, faultToleranceLevel int, workersAddresses []string) *Task {

	output := new(Task)

	(*output).workersReplyChannel = make(chan workerReply)
	(*output).faultToleranceLevel = faultToleranceLevel
	(*output).currentWorkerID = 0
	(*output).workersAddresses = workersAddresses
	(*output).repliesRegistry = WorkersResponsesRegistry.New((*output).faultToleranceLevel + 1)
	(*output).inputSplit = inputSplit

	return output
}

/**
 * This function is used to execute a 'Byzantine Fault Tolerant' Map-task.
 */
func (obj *Task) Execute() (string, []string) {

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

func (obj *Task) startListeningWorkersReplies() {

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

func executeSingleMapTaskReplica(middleInput TransientData, address string, reply chan workerReply) {

	var workerInput ReplicaInput
	var workerOutput ReplicaOutput

	worker, err := rpc.Dial("tcp", address)
	utility.CheckError(err)

	workerInput.Data = middleInput

	err = worker.Call("SingleMapReduceReplica.Execute", &workerInput, &workerOutput)

	if err == nil {

		output := new(workerReply)
		(*output).workerAddress = address
		(*output).digest = workerOutput.Digest

		reply <- *output
	}
}
