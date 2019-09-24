package primary

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount"
	"encoding/gob"
	"fmt"
)

func StartWork() {

	var err error
	var allPendingClientRequestGuid []string

	process.GetPropertyAsString(property.InternetAddress)

	address := fmt.Sprintf("%s:%d", process.GetPropertyAsString(property.InternetAddress), aftmapreduce.WordCountRequestRPCBasePort+process.GetPropertyAsInteger(property.NodeID))

	err = (*process.GetSystemCoordinator()).WaitUntilLeader(address)
	utility.CheckError(err)

	(*process.GetLogger()).PrintInfoLevelLabeledMessage("Initialization", "I'm leader")

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	allPendingClientRequestGuid, err = (*process.GetSystemCoordinator()).GetAllPendingClientRequestGuid()
	utility.CheckError(err)

	for _, request := range allPendingClientRequestGuid {
		go wordcount.JobStart(request)
	}

	aftmapreduce.StartAcceptingRPCRequest(&wordcount.Request{}, process.GetPropertyAsString(property.WordCountRequestRPCFullAddress))
}
