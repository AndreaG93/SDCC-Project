package Algorithm

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/DataStructures/WorkerResponseRegistry"
	"SDCC-Project-WorkerNode/utility"
	"net/rpc"
)

var faultToleranceLevel = 5
var maximumNumberOfFaultyReplicas = 5

func dddd() {

}

func startBFTMapPhase(inputData string) {

	var requiredNumberOfMatchesReached chan bool
	var registry *WorkersResponseRegistry.WorkersResponseRegistry

	requiredNumberOfMatchesReached = make(chan bool)

	registry = WorkersResponseRegistry.New(faultToleranceLevel, requiredNumberOfMatchesReached)

	StartMapTaskReplicas()

	for replica := 0; replica < maximumNumberOfFaultyReplicas; replica++ {

		go func() {

		}()
	}

	<-requiredNumberOfMatchesReached

}

func startSingleMapTaskReplica(workerAddress string, inputData string) string {

	var mapTaskInput wordcount.MapInput
	var mapTaskOutput wordcount.MapOutput

	worker, err := rpc.Dial("tcp", workerAddress)
	utility.CheckError(err)

	mapTaskInput.Input = inputData
	mapTaskInput.MapCardinality = 5

	err = worker.Call("Map.Execute", &mapTaskInput, &mapTaskOutput)
	utility.CheckError(err)

	return mapTaskOutput.Digest
}
