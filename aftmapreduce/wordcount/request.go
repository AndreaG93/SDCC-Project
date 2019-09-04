package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"io/ioutil"
	"os"
)

type Request struct {
}

type RequestInput struct {
	FileContent      string
	SourceFileDigest string
}

type RequestOutput struct {
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	rawData := []byte(input.FileContent)

	file, err := ioutil.TempFile("", "")
	utility.CheckError(err)

	_, err = file.Write(rawData)
	utility.CheckError(err)

	utility.CheckError(file.Sync())

	_, err = file.Seek(0, 0)
	utility.CheckError(err)

	node.GetDataRegistry().Set(input.SourceFileDigest, rawData)
	node.GetAmazonS3Client().Upload(file, input.SourceFileDigest)

	utility.CheckError(file.Close())
	utility.CheckError(os.Remove(file.Name()))

	go ManageRequest(NewClientRequest(input.SourceFileDigest))

	return nil
}

func ManageRequest(clientRequest *ClientRequest) {

	var mapTaskOutput []*AFTMapTaskOutput
	var reduceTaskOutput []*AFTReduceTaskOutput
	var localityAwarenessData map[int]int
	var rawData []byte

	for {

		currentClientRequestStatus := clientRequest.getStatus()

		node.SetProperty(property.MapCardinality, node.GetZookeeperClient().GetGroupAmount())
		node.GetLogger().PrintInfoStartingTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))

		switch currentClientRequestStatus {
		case InitialStatus:

			splits := getSplits(clientRequest.digest, node.GetPropertyAsInteger(property.MapCardinality))
			mapTaskOutput = mapTask(splits)

			rawData = utility.Encode(mapTaskOutput)
			(*clientRequest).CheckPoint(AfterMapStatus, rawData, nil)

		case AfterMapStatus:

			if mapTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache1()
				utility.Decode(rawData, &mapTaskOutput)
			}

			node.GetAmazonS3Client().Delete((*clientRequest).digest)

			localityAwarenessData = getLocalityAwareReduceTaskMappedToNodeGroupId(mapTaskOutput)
			localityAwareShuffleTask(mapTaskOutput, localityAwarenessData)

			rawData = utility.Encode(localityAwarenessData)

			(*clientRequest).CheckPoint(AfterLocalityAwareShuffle, nil, rawData)

		case AfterLocalityAwareShuffle:

			if localityAwarenessData == nil || mapTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache1()
				utility.Decode(rawData, &mapTaskOutput)

				rawData = (*clientRequest).GetDataFromCache2()
				utility.Decode(rawData, &localityAwarenessData)
			}

			reduceTaskOutput = reduceTask(mapTaskOutput, localityAwarenessData)

			rawData = utility.Encode(reduceTaskOutput)

			(*clientRequest).CheckPoint(AfterReduce, rawData, nil)

		case AfterReduce:

			if reduceTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache1()
				utility.Decode(rawData, &reduceTaskOutput)
			}

			dataArray := retrieveTask(reduceTaskOutput)

			finalOutput := computeFinalOutputTask(dataArray)
			finalRawData := finalOutput.Serialize()

			node.GetZookeeperClient().SetZNodeData((*clientRequest).GetCompleteRequestZNodePath(), finalRawData)

			(*clientRequest).CheckPoint(Complete, nil, nil)

		case Complete:
			clientRequest.DeletePendingRequest()
			return
		}
		node.GetLogger().PrintInfoCompleteTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))
	}
}
