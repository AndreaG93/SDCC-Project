package aft

import (
	"SDCC-Project/aftmapreduce/registry"
	"SDCC-Project/aftmapreduce/wordcount"
	"SDCC-Project/utility"
	"net/rpc"
)

type MapTaskOutput struct {
	IdGroup                  int
	ReplayDigest             string
	NodeIdsWithCorrectResult []int
	MappedDataSizes          map[int]int
}

type reply struct {
	id              int
	digest          string
	mappedDataSizes map[int]int
}

type MapTask struct {
	workersReplyChannel chan workerReply
	faultToleranceLevel int
	currentWorkerID     int
	workersAddresses    []string
	registry            *registry.MapReplies
	split               string
}

func NewMapTask(split string, workersAddresses []string) *ReplicatedMapTask {

	output := new(ReplicatedMapTask)

	(*output).workersReplyChannel = make(chan workerReply)
	(*output).faultToleranceLevel = faultToleranceLevel
	(*output).currentWorkerID = 0
	(*output).workersAddresses = workersAddresses
	(*output).registry = registry.NewMapReply(faultToleranceLevel + 1)
	(*output).split = split

	return output
}

func (obj *ReplicatedMapTask) Execute() (string, []string) {

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

func (obj *ReplicatedMapTask) startListeningWorkersReplies() {

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

func executeSingleMapTaskReplica(split string, address string, reply chan workerReply) {

	var workerInput wordcount.MapInput
	var workerOutput wordcount.MapOutput

	worker, err := rpc.Dial("tcp", address)
	utility.CheckError(err)

	workerInput.text = middleInput

	err = worker.Call("Replica.Execute", &workerInput, &workerOutput)

	if err == nil {

		output := new(workerReply)
		(*output).workerAddress = address
		(*output).digest = workerOutput.Digest

		reply <- *output
	}
}
