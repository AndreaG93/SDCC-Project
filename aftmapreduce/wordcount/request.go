package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/utility"
	"fmt"
)

type Request struct {
}

type RequestInput struct {
	SourceFileDigest string
}

type RequestOutput struct {
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	go ManageRequest(NewClientRequest(input.SourceFileDigest))

	return nil
}

func ManageRequest(clientRequest *ClientRequest) {

	var mapTaskOutput []*AFTMapTaskOutput
	var reduceTaskOutput []*AFTReduceTaskOutput
	var localityAwarenessData map[int]int

	for {

		currentClientRequestStatus := clientRequest.getStatus()

		node.SetProperty(property.MapCardinality, node.GetZookeeperClient().GetGroupAmount())
		node.GetLogger().PrintInfoStartingTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))

		switch currentClientRequestStatus {
		case InitialStatus:

			splits := getSplits(clientRequest.digest, node.GetPropertyAsInteger(property.MapCardinality))
			mapTaskOutput = mapTask(splits)

			(*clientRequest).CheckPoint(AfterMapStatus, nil)

		case AfterMapStatus:

			node.GetAmazonS3Client().Delete((*clientRequest).digest)

			localityAwarenessData = getLocalityAwareReduceTaskMappedToNodeGroupId(mapTaskOutput)
			localityAwareShuffleTask(mapTaskOutput, localityAwarenessData)

			(*clientRequest).CheckPoint(AfterLocalityAwareShuffle, nil)

		case AfterLocalityAwareShuffle:

			reduceTaskOutput = reduceTask(mapTaskOutput, localityAwarenessData)

			(*clientRequest).CheckPoint(AfterLocalityAwareShuffle, nil)

		case AfterReduce:
			dataArray := retrieveTask(reduceTaskOutput)

			finalOutput := computeFinalOutputTask(dataArray)
			finalRawData, err := finalOutput.Serialize()
			utility.CheckError(err)

			node.GetZookeeperClient().SetZNodeData((*clientRequest).GetCompleteRequestZNodePath(), finalRawData)

			(*clientRequest).CheckPoint(Complete, nil)

		case Complete:
			return
		}
		node.GetLogger().PrintInfoCompleteTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))
	}
}
