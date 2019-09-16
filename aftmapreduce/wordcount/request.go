package wordcount

import (
	"SDCC-Project/aftmapreduce/cloud/zookeeper"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/utility"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	// Request's types
	AcceptanceJobRequestType        = uint8(0)
	UploadPreSignedURLRequestType   = uint8(1)
	DownloadPreSignedURLRequestType = uint8(2)

	// Request's status
)

type Request struct {
}

type RequestInput struct {
	Type             uint8
	SourceFileDigest string
}

type RequestOutput struct {
	PreSignedURL string
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	switch input.Type {
	case UploadPreSignedURLRequestType:
		output.PreSignedURL = node.GetAmazonS3Client().GetPreSignedURLForUploadTask(input.SourceFileDigest)
		return nil
	case DownloadPreSignedURLRequestType:
		output.PreSignedURL = node.GetAmazonS3Client().GetPreSignedURLForDownloadTask(input.SourceFileDigest)
		return nil
	case AcceptanceJobRequestType:

		if isRequestAlreadyAccepted(input.SourceFileDigest) {
			return nil
		} else {
			myClientRequest := zookeeper.NewClientRequest(input.SourceFileDigest)
			go ManageClientRequest(input.SourceFileDigest)
		}

	default:
		return errors.New("request type not recognized")
	}

	return nil
}

func ManageClientRequest(guid string) {

	var mapTaskOutput []*AFTMapTaskOutput
	var reduceTaskOutput []*AFTReduceTaskOutput
	var localityAwarenessData map[int]int
	var rawData []byte

	for {

		currentClientRequestStatus := clientRequest.getStatus()

		node.SetProperty(property.MapCardinality, node.GetZookeeperClient().GetGroupAmount())
		node.GetLogger().PrintInfoStartingTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))

		switch currentClientRequestStatus {
		case zookeeper.InitialStatus:

			splits := getSplits(clientRequest.digest, node.GetPropertyAsInteger(property.MapCardinality))
			mapTaskOutput = mapTask(splits)

			rawData = utility.Encode(mapTaskOutput)
			(*clientRequest).CheckPoint(zookeeper.AfterMap, rawData)

		case zookeeper.AfterMap:

			if mapTaskOutput == nil {
				rawData = (*clientRequest).GetDataFromCache()
				utility.Decode(rawData, &mapTaskOutput)
			}

			node.GetAmazonS3Client().Delete((*clientRequest).digest)

			localityAwarenessData = getLocalityAwareReduceTaskMappedToNodeGroupId(mapTaskOutput)
			localityAwareShuffleTask(mapTaskOutput, localityAwarenessData)

			reduceTaskOutput = reduceTask(mapTaskOutput, localityAwarenessData)

			rawData = utility.Encode(reduceTaskOutput)

			(*clientRequest).CheckPoint(zookeeper.AfterReduce, rawData)

		case zookeeper.AfterReduce:

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

			(*clientRequest).CheckPoint(zookeeper.Complete, nil)

		case zookeeper.Complete:
			clientRequest.DeletePendingRequest()
			return
		}
		node.GetLogger().PrintInfoCompleteTaskMessage(fmt.Sprintf("%s for %s", currentClientRequestStatus, (*clientRequest).digest))
	}
}
