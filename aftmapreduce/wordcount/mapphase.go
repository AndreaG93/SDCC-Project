package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"sync"
)

func startMapPhase(guid string) error {

	var syncWaitGroup sync.WaitGroup

	if splits, err := getSplits(guid, (*process.GetMembershipRegister()).GetGroupAmount()); err != nil {
		return err
	} else {

		output := make([]*AFTMapTaskOutput, (*process.GetMembershipRegister()).GetGroupAmount())
		firstRepliesChannel := make(chan interface{})
		isSpeculativeExecutionSucceeded := make(chan bool)
		mapPhaseOutputChannel := make(chan []*AFTMapTaskOutput)

		go startMapPhaseWithSpeculativeExecution(guid, firstRepliesChannel, isSpeculativeExecutionSucceeded, mapPhaseOutputChannel)

		for workerProcessGroupID, split := range splits {
			syncWaitGroup.Add(1)
			go startAFTMapTask(split, workerProcessGroupID, &output[workerProcessGroupID], &syncWaitGroup, firstRepliesChannel)
		}

		syncWaitGroup.Wait()

		mapPhaseOutputChannel <- output

		if !<-isSpeculativeExecutionSucceeded {

			process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Map-Phase with speculative execution FAILED for Request GUID: %s", guid))

			if rawOutput, err := utility.Encoding(output); err != nil {
				return err
			} else {
				if err := (*process.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, mapPhaseComplete, rawOutput); err != nil {
					return err
				} else {
					process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Map-Phase with Standard execution Complete for Request GUID: %s", guid))
				}
			}
		}

		close(mapPhaseOutputChannel)
		close(isSpeculativeExecutionSucceeded)

		return nil
	}
}

func startAFTMapTask(split string, workerProcessGroupID int, output **AFTMapTaskOutput, syncWaitGroup *sync.WaitGroup, firstRepliesChannel chan interface{}) {
	*output = Execute(NewAFTMapTask(split, workerProcessGroupID, firstRepliesChannel)).(*AFTMapTaskOutput)
	syncWaitGroup.Done()
}

func startMapPhaseWithSpeculativeExecution(guid string, firstRepliesChannel chan interface{}, predictionSuccessfulChannel chan bool, mapPhaseOutputChannel chan []*AFTMapTaskOutput) {

	reducePhaseOutputRawDataChannel := make(chan []byte)
	remainingRepliesUntilCompletion := (*process.GetMembershipRegister()).GetGroupAmount()
	output := buildAFTMapTaskOutputsArray()

	for reply := range firstRepliesChannel {

		remainingRepliesUntilCompletion--

		firstMapOutputReply := reply.(*MapOutput)
		workerProcessGroupID := (*firstMapOutputReply).IdGroup

		output[workerProcessGroupID].IdGroup = workerProcessGroupID
		output[workerProcessGroupID].ReplayDigest = (*firstMapOutputReply).ReplayDigest
		output[workerProcessGroupID].NodeIdsWithCorrectResult = []int{(*firstMapOutputReply).IdNode}
		output[workerProcessGroupID].MappedDataSizes = (*firstMapOutputReply).MappedDataSizes

		if remainingRepliesUntilCompletion == 0 {
			break
		}
	}

	close(firstRepliesChannel)
	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Map-Phase with speculative execution complete (but not yet confimed) for Request GUID: %s", guid))

	go startReducePhaseAfterMapPhaseSpeculativeExecution(guid, output, reducePhaseOutputRawDataChannel)

	mapPhaseOutputWithStandardAlgorithm := <-mapPhaseOutputChannel

	if !areMapPhaseOutputEqual(mapPhaseOutputWithStandardAlgorithm, output) {
		predictionSuccessfulChannel <- false
		<-reducePhaseOutputRawDataChannel
	} else {

		rawReducePhaseOutput := <-reducePhaseOutputRawDataChannel

		if rawReducePhaseOutput != nil {

			if err := (*process.GetSystemCoordinator()).UpdateClientRequestStatusBackup(guid, reducePhaseComplete, rawReducePhaseOutput); err != nil {
				predictionSuccessfulChannel <- false
			} else {
				process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Map-Phase with speculative execution FULLY complete for Request GUID: %s", guid))
				predictionSuccessfulChannel <- true
			}

		} else {
			predictionSuccessfulChannel <- false
		}
	}

	close(reducePhaseOutputRawDataChannel)
}

func startReducePhaseAfterMapPhaseSpeculativeExecution(guid string, mapPhaseOutput []*AFTMapTaskOutput, reducePhaseOutputRawDataChannel chan []byte) {

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Reduce-Phase AFTER speculative execution STARTED for Request GUID: %s", guid))

	localityAwareReduceTaskSchedule := getLocalityAwareReduceTaskSchedule(mapPhaseOutput)
	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("The Locality-Aware AFT-Reduce Task Schedule (For Speculative Execution) for Request %s is :: %d", guid, localityAwareReduceTaskSchedule))

	startShuffleTask(mapPhaseOutput, localityAwareReduceTaskSchedule)

	output := reduceTask(mapPhaseOutput, localityAwareReduceTaskSchedule)

	if rawOutput, err := utility.Encoding(output); err != nil {
		reducePhaseOutputRawDataChannel <- nil
	} else {
		reducePhaseOutputRawDataChannel <- rawOutput
	}

	process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Reduce-Phase AFTER speculative execution FINISHED for Request GUID: %s", guid))
}

func areMapPhaseOutputEqual(a []*AFTMapTaskOutput, b []*AFTMapTaskOutput) bool {

	output := true

	for x := range a {
		if b[x].IdGroup != a[x].IdGroup {
			output = false
			break
		}
		if b[x].ReplayDigest != a[x].ReplayDigest {
			output = false
			break
		}
	}

	return output
}

func buildAFTMapTaskOutputsArray() []*AFTMapTaskOutput {

	output := make([]*AFTMapTaskOutput, (*process.GetMembershipRegister()).GetGroupAmount())
	for x := range output {
		output[x] = new(AFTMapTaskOutput)
	}

	return output
}
