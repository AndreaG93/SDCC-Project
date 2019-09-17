package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/utility"
	"errors"
	"fmt"
)

const (
	requestManagementTask = "REQUEST MANAGEMENT"

	// Request's types
	AcceptanceJobRequestType        = uint8(0)
	UploadPreSignedURLRequestType   = uint8(1)
	DownloadPreSignedURLRequestType = uint8(2)

	// Request's status
	accepted       = uint8(0)
	mapComplete    = uint8(1)
	reduceComplete = uint8(2)
	complete       = uint8(3)
)

type Request struct {
}

type RequestInput struct {
	Type             uint8
	SourceFileDigest string
}

type RequestOutput struct {
	Url string
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	var err error
	var isRequestAlreadyAccepted bool

	switch input.Type {
	case UploadPreSignedURLRequestType:
		output.Url, err = (*node.GetKeyValueStorageService()).RetrieveURLForPutOperation(input.SourceFileDigest)
	case DownloadPreSignedURLRequestType:
		output.Url, err = (*node.GetKeyValueStorageService()).RetrieveURLForGetOperation(input.SourceFileDigest)
	case AcceptanceJobRequestType:

		isRequestAlreadyAccepted, err = (*node.GetSystemCoordinator()).ClientRequestExist(input.SourceFileDigest)
		if err != nil || isRequestAlreadyAccepted {
			break
		} else {
			if err = (*node.GetSystemCoordinator()).RegisterClientRequest(input.SourceFileDigest, accepted); err == nil {
				go ManageClientRequest(input.SourceFileDigest)
			}
		}
	default:
		return errors.New("request type not recognized")
	}

	return err
}

func ManageClientRequest(guid string) {

	var status uint8
	var data []byte
	var err error

	var AFTMapTaskOutput []*AFTMapTaskOutput
	var AFTReduceTaskOutput []*AFTReduceTaskOutput

	status, data, err = (*node.GetSystemCoordinator()).GetClientRequestInformation(guid)
	utility.CheckError(err)

	node.GetLogger().PrintInfoTaskMessage(requestManagementTask, fmt.Sprintf("Request: %s -- Status %d", guid, status))

	for {
		switch status {
		case accepted:

			if AFTMapTaskOutput, err = startAFTMapTask(guid); err == nil {
				status = mapComplete
			}

		case mapComplete:

			if AFTMapTaskOutput == nil {
				utility.Decode(data, &AFTMapTaskOutput)
			}

			if AFTReduceTaskOutput, err = startAFTReduceTask(guid, AFTMapTaskOutput); err == nil {
				status = reduceComplete
			}

		case reduceComplete:

			if AFTReduceTaskOutput == nil {
				utility.Decode(data, &AFTReduceTaskOutput)
			}

			if err = startCollectTask(guid, AFTReduceTaskOutput); err == nil {
				status = complete
			}

		case complete:
			if err = (*node.GetSystemCoordinator()).DeletePendingRequest(guid); err == nil {
				return
			}
		}
	}
}

func startAFTMapTask(guid string) ([]*AFTMapTaskOutput, error) {

	splits := getSplits(guid, (*node.GetMembershipRegister()).GetGroupAmount())
	output := mapTask(splits)

	return output, (*node.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, mapComplete, utility.Encode(output))
}

func startAFTReduceTask(guid string, AFTMapTaskOutput []*AFTMapTaskOutput) ([]*AFTReduceTaskOutput, error) {

	localityAwarenessData := getLocalityAwareReduceTaskMappedToNodeGroupId(AFTMapTaskOutput)
	localityAwareShuffleTask(AFTMapTaskOutput, localityAwarenessData)

	output := reduceTask(AFTMapTaskOutput, localityAwarenessData)

	return output, (*node.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, reduceComplete, utility.Encode(output))
}

func startCollectTask(guid string, AFTReduceTaskOutput []*AFTReduceTaskOutput) error {

	var err error

	dataArray := retrieveTask(AFTReduceTaskOutput)

	output := computeFinalOutputTask(dataArray)
	outputSerialized := output.Serialize()
	outputGUID := utility.GenerateDigestUsingSHA512(outputSerialized)

	if err = (*node.GetKeyValueStorageService()).Put(outputGUID, outputSerialized); err == nil {
		if err = (*node.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, complete, nil); err == nil {
			err = (*node.GetSystemCoordinator()).RegisterClientRequestAsComplete(guid, outputGUID)
		}
	}

	return err
}
