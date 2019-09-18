package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/utility"
	"errors"
	"fmt"
	"sync"
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
		output.Url, err = (*node.GetStorageKeyValueRegister()).RetrieveURLForPutOperation(input.SourceFileDigest)
	case DownloadPreSignedURLRequestType:
		output.Url, err = (*node.GetStorageKeyValueRegister()).RetrieveURLForGetOperation(input.SourceFileDigest)
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

	var myAFTMapTaskOutput []*AFTMapTaskOutput
	var myAFTReduceTaskOutput []*AFTReduceTaskOutput

	predictionSuccessfulChannel := make(chan bool)
	AFTMapTaskOutputComputedWithStandardAlgorithm := make(chan []*AFTMapTaskOutput)

	for {

		status, data, err = (*node.GetSystemCoordinator()).GetClientRequestInformation(guid)
		utility.CheckError(err)
		node.GetLogger().PrintInfoTaskMessage(requestManagementTask, fmt.Sprintf("Request: %s -- Status %d", guid, status))

		switch status {
		case accepted:

			if myAFTMapTaskOutput, err = startAFTMapTask(guid, AFTMapTaskOutputComputedWithStandardAlgorithm, predictionSuccessfulChannel); err == nil {

				predictedResultIsCorrect := <-predictionSuccessfulChannel
				if predictedResultIsCorrect {
					status = reduceComplete
				} else {
					status = mapComplete
				}
			}

		case mapComplete:

			if myAFTMapTaskOutput == nil {
				utility.Decode(data, &myAFTMapTaskOutput)
			}

			if myAFTReduceTaskOutput, err = startAFTReduceTask(guid, myAFTMapTaskOutput); err == nil {
				status = reduceComplete
			}

		case reduceComplete:

			if myAFTReduceTaskOutput == nil {
				utility.Decode(data, &myAFTReduceTaskOutput)
			}

			if err = startCollectTask(guid, myAFTReduceTaskOutput); err == nil {
				status = complete
			}

		case complete:
			if err = (*node.GetSystemCoordinator()).DeletePendingRequest(guid); err == nil {
				return
			}
		}
	}
}

func startAFTMapTask(guid string, AFTMapTaskOutputComputedWithStandardAlgorithm chan []*AFTMapTaskOutput, predictionSuccessfulChannel chan bool) ([]*AFTMapTaskOutput, error) {

	var mapWaitGroup sync.WaitGroup

	if splits, err := getSplits(guid, (*node.GetMembershipRegister()).GetGroupAmount()); err != nil {
		return nil, err
	} else {

		output := make([]*AFTMapTaskOutput, len(splits))
		channels := getChannelsUsedForFirstReplyPredictedAsCorrect(guid, len(splits), AFTMapTaskOutputComputedWithStandardAlgorithm, predictionSuccessfulChannel)

		for index := range splits {

			mapWaitGroup.Add(1)
			go func(mySplit string, myChannel chan interface{}, myGroupId int) {

				output[myGroupId] = Execute(NewMapTask(mySplit, myGroupId, myChannel)).(*AFTMapTaskOutput)
				mapWaitGroup.Done()

			}(splits[index], channels[index], index)
		}

		mapWaitGroup.Wait()
		node.GetLogger().PrintInfoCompleteTaskMessage(MapTaskName)

		AFTMapTaskOutputComputedWithStandardAlgorithm <- output

		return output, (*node.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, mapComplete, utility.Encode(output))
	}
}

func getChannelsUsedForFirstReplyPredictedAsCorrect(guid string, splitsAmount int, AFTMapTaskOutputComputedWithStandardAlgorithm chan []*AFTMapTaskOutput, predictionSuccessfulChannel chan bool) []chan interface{} {

	var mapWaitGroup sync.WaitGroup
	isPredictionCorrect := true

	output := make([]chan interface{}, splitsAmount)
	for index := range output {
		output[index] = make(chan interface{})
	}

	go func(myChannel []chan interface{}) {

		myAFTMapTaskOutput := make([]*AFTMapTaskOutput, splitsAmount)

		for index := range myChannel {

			mapWaitGroup.Add(1)
			go func(myIndex int, myChannel chan interface{}) {

				reply := <-myChannel

				myAFTMapTaskOutput[myIndex] = new(AFTMapTaskOutput)
				myAFTMapTaskOutput[myIndex].IdGroup = myIndex
				myAFTMapTaskOutput[myIndex].ReplayDigest = reply.(*MapOutput).ReplayDigest
				myAFTMapTaskOutput[myIndex].NodeIdsWithCorrectResult = []int{reply.(*MapOutput).IdNode}
				myAFTMapTaskOutput[myIndex].MappedDataSizes = reply.(*MapOutput).MappedDataSizes

				mapWaitGroup.Done()

			}(index, myChannel[index])
		}
		mapWaitGroup.Wait()

		localityAwarenessData := getLocalityAwareReduceTaskMappedToNodeGroupId(myAFTMapTaskOutput)
		localityAwareShuffleTask(myAFTMapTaskOutput, localityAwarenessData)

		output := reduceTask(myAFTMapTaskOutput, localityAwarenessData)

		mapTaskFromStandardAlgorithm := <-AFTMapTaskOutputComputedWithStandardAlgorithm

		for x := range myAFTMapTaskOutput {
			if myAFTMapTaskOutput[x].IdGroup != mapTaskFromStandardAlgorithm[x].IdGroup {
				isPredictionCorrect = false
				break
			}
			if myAFTMapTaskOutput[x].ReplayDigest != mapTaskFromStandardAlgorithm[x].ReplayDigest {
				isPredictionCorrect = false
				break
			}
		}

		if isPredictionCorrect {
			utility.CheckError((*node.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, reduceComplete, utility.Encode(output)))
			predictionSuccessfulChannel <- true
		} else {
			predictionSuccessfulChannel <- false
		}

	}(output)

	return output
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

	if err = (*node.GetStorageKeyValueRegister()).Put(outputGUID, outputSerialized); err == nil {
		if err = (*node.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, complete, nil); err == nil {
			err = (*node.GetSystemCoordinator()).RegisterClientRequestAsComplete(guid, outputGUID)
		}
	}

	return err
}
