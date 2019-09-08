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
	FileContent                    string
	RequestPreSignedURLForUpload   bool
	RequestPreSignedURLForDownload bool
	SourceFileDigest               string
}

type RequestOutput struct {
	PreSignedURL string
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	if input.RequestPreSignedURLForUpload {
		output.PreSignedURL = node.GetAmazonS3Client().GetPreSignedURL(input.SourceFileDigest)
		return nil
	}

	if input.RequestPreSignedURLForDownload {
		output.PreSignedURL = node.GetAmazonS3Client().GetPreSignedURLForDownloadOperation(input.SourceFileDigest)
		return nil
	}

	if CheckDuplicatedClientRequest(input.SourceFileDigest) {
		return nil
	}

	myClientRequest := NewClientRequest(input.SourceFileDigest)
	go ManageRequest(myClientRequest)

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
			(*clientRequest).CheckPoint(AfterMap, rawData)

		case AfterMap:

			if mapTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache()
				utility.Decode(rawData, &mapTaskOutput)
			}

			node.GetAmazonS3Client().Delete((*clientRequest).digest)

			localityAwarenessData = getLocalityAwareReduceTaskMappedToNodeGroupId(mapTaskOutput)
			localityAwareShuffleTask(mapTaskOutput, localityAwarenessData)

			reduceTaskOutput = reduceTask(mapTaskOutput, localityAwarenessData)

			rawData = utility.Encode(reduceTaskOutput)

			(*clientRequest).CheckPoint(AfterReduce, rawData)

		case AfterReduce:

			if reduceTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache()
				utility.Decode(rawData, &reduceTaskOutput)
			}

			dataArray := retrieveTask(reduceTaskOutput)

			finalOutput := computeFinalOutputTask(dataArray)
			finalRawData := finalOutput.Serialize()

			finalOutputDigestData := utility.GenerateDigestUsingSHA512(finalRawData)

			output, err := ioutil.TempFile(os.TempDir(), finalOutputDigestData)
			utility.CheckError(err)

			_, err = output.Write(finalRawData)
			utility.CheckError(err)
			utility.CheckError(output.Sync())
			_, err = output.Seek(0, 0)
			utility.CheckError(err)

			node.GetAmazonS3Client().Upload(output, finalOutputDigestData)

			node.GetZookeeperClient().SetZNodeData((*clientRequest).GetCompleteRequestZNodePath(), []byte(finalOutputDigestData))

			(*clientRequest).CheckPoint(Complete, nil)

		case Complete:
			clientRequest.DeletePendingRequest()
			return
		}
		node.GetLogger().PrintInfoCompleteTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))
	}
}
