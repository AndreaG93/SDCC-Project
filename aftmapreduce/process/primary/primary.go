package primary

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount"
	"encoding/gob"
)

func StartWork() {

	var err error
	var allPendingClientRequestGuid []string

	err = (*process.GetSystemCoordinator()).WaitUntilLeader(process.GetPropertyAsString(property.WordCountRequestRPCFullAddress))
	utility.CheckError(err)

	(*process.GetLogger()).PrintInfoLevelLabeledMessage("Initialization", "I'm leader")

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	allPendingClientRequestGuid, err = (*process.GetSystemCoordinator()).GetAllPendingClientRequestGuid()
	utility.CheckError(err)

	for _, request := range allPendingClientRequestGuid {
		go wordcount.ManageClientRequest(request)
	}

	aftmapreduce.StartAcceptingRPCRequest(&wordcount.Request{}, process.GetPropertyAsString(property.WordCountRequestRPCFullAddress))
}
