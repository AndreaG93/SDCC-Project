package worker

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount"
)

func StartWork() {

	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Reduce{}, process.GetPropertyAsString(property.WordCountReduceRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Receive{}, process.GetPropertyAsString(property.WordCountReceiveRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Send{}, process.GetPropertyAsString(property.WordCountSendRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Retrieve{}, process.GetPropertyAsString(property.WordCountRetrieveRPCFullAddress))

	err := (*process.GetMembershipCoordinator()).RegisterNewWorkerProcess(process.GetPropertyAsInteger(property.NodeID), process.GetPropertyAsInteger(property.NodeGroupID), process.GetPropertyAsString(property.InternetAddress))
	utility.CheckError(err)

	aftmapreduce.StartAcceptingRPCRequest(&wordcount.Map{}, process.GetPropertyAsString(property.WordCountMapRPCFullAddress))
}
