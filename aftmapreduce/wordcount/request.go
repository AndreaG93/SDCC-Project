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
	var rawData []byte
	var err error

	for {

		reduceTaskOutput = nil
		mapTaskOutput = nil
		localityAwarenessData = nil

		currentClientRequestStatus := clientRequest.getStatus()

		node.SetProperty(property.MapCardinality, node.GetZookeeperClient().GetGroupAmount())
		node.GetLogger().PrintInfoStartingTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))

		switch currentClientRequestStatus {
		case InitialStatus:

			splits := getSplits(clientRequest.digest, node.GetPropertyAsInteger(property.MapCardinality))
			mapTaskOutput = mapTask(splits)

			rawData, err = utility.Encode(mapTaskOutput)
			utility.CheckError(err)

			(*clientRequest).CheckPoint(AfterMapStatus, rawData, nil)

		case AfterMapStatus:

			if mapTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache1()
				utility.CheckError(utility.Decode(rawData, &mapTaskOutput))
			}

			node.GetAmazonS3Client().Delete((*clientRequest).digest)

			localityAwarenessData = getLocalityAwareReduceTaskMappedToNodeGroupId(mapTaskOutput)
			localityAwareShuffleTask(mapTaskOutput, localityAwarenessData)

			rawData, err = utility.Encode(localityAwarenessData)
			utility.CheckError(err)

			(*clientRequest).CheckPoint(AfterLocalityAwareShuffle, nil, rawData)

		case AfterLocalityAwareShuffle:

			if localityAwarenessData == nil || mapTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache1()
				utility.CheckError(utility.Decode(rawData, &mapTaskOutput))

				rawData = (*clientRequest).GetDataFromCache2()
				utility.CheckError(utility.Decode(rawData, &localityAwarenessData))
			}

			reduceTaskOutput = reduceTask(mapTaskOutput, localityAwarenessData)

			rawData, err = utility.Encode(reduceTaskOutput)
			utility.CheckError(err)

			(*clientRequest).CheckPoint(AfterReduce, rawData, nil)

		case AfterReduce:

			if reduceTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache1()
				utility.CheckError(utility.Decode(rawData, &reduceTaskOutput))
			}

			dataArray := retrieveTask(reduceTaskOutput)

			finalOutput := computeFinalOutputTask(dataArray)
			finalRawData, err := finalOutput.Serialize()
			utility.CheckError(err)

			node.GetZookeeperClient().SetZNodeData((*clientRequest).GetCompleteRequestZNodePath(), finalRawData)

			(*clientRequest).CheckPoint(Complete, nil, nil)

		case Complete:
			clientRequest.DeletePendingRequest()
			return
		}
		node.GetLogger().PrintInfoCompleteTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))
	}
}
