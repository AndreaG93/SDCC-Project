package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
)

const (
	acceptancePhaseComplete = uint8(0)
	mapPhaseComplete        = uint8(1)
	reducePhaseComplete     = uint8(2)
	collectPhaseComplete    = uint8(3)
)

func JobStart(guid string) {

	var status uint8
	var data []byte
	var err error

	var mapPhaseOutput []*AFTMapTaskOutput
	var reducePhaseOutput []*AFTReduceTaskOutput

	for {

		if status, data, err = (*process.GetSystemCoordinator()).GetClientRequestInformation(guid); err != nil {
			process.GetLogger().PrintErrorLevelMessage(err.Error())
			continue
		}

		switch status {
		case acceptancePhaseComplete:

			process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Start Map-Phase for Request GUID: %s", guid))

			if err = startMapPhase(guid); err != nil {
				process.GetLogger().PrintErrorLevelMessage(err.Error())
			}

		case mapPhaseComplete:

			process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Start Reduce-Phase for Request GUID: %s", guid))

			if mapPhaseOutput == nil {
				if err = utility.Decoding(data, &mapPhaseOutput); err != nil {
					process.GetLogger().PrintErrorLevelMessage(err.Error())
					continue
				}
			}

			if reducePhaseOutput, err = startReducePhase(guid, mapPhaseOutput); err == nil {
				status = reducePhaseComplete
			}

		case reducePhaseComplete:

			process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Start Collect-Phase for Request GUID: %s", guid))

			if reducePhaseOutput == nil {
				if err = utility.Decoding(data, &reducePhaseOutput); err != nil {
					process.GetLogger().PrintErrorLevelMessage(err.Error())
					continue
				}
			}

			if err = startCollectPhase(guid, reducePhaseOutput); err == nil {
				status = collectPhaseComplete
			}

		case collectPhaseComplete:

			if err = (*process.GetSystemCoordinator()).DeletePendingRequest(guid); err == nil {

				process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Job complete for for Request GUID: %s", guid))
				return
			}
		}
	}
}
