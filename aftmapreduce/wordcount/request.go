package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
)

type Request struct {
}

type RequestInput struct {
	SourceFileDigest string
}

type RequestOutput struct {
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	go manageRequest(input.SourceFileDigest)
	return nil
}

func manageRequest(digest string) {

	node.SetProperty(property.MapCardinality, node.GetZookeeperClient().GetGroupAmount())

	splits := getSplits(digest, node.GetPropertyAsInteger(property.MapCardinality))

	mapTaskOutput := mapTask(splits)
	localityAwarenessData := getLocalityAwareReduceTaskMappedToNodeGroupId(mapTaskOutput)
	reduceTaskOutput := localityAwareShuffleAndReduceTask(mapTaskOutput, localityAwarenessData)

	dataArray := retrieveTask(reduceTaskOutput)
	finalData := computeFinalOutputTask(dataArray)

	finalData.Print()

}
