package primary

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount"
	"encoding/gob"
	"fmt"
)

func Initialize(id int, internetAddress string, zookeeperAddresses []string) {

	utility.CheckError(node.InitializePrimary(uint(id), zookeeperAddresses))

	node.SetProperty(property.NodeID, id)
	node.SetProperty(property.NodeType, "Primary")
	node.SetProperty(property.InternetAddress, internetAddress)
	node.SetProperty(property.WordCountRequestRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountRequestRPCBasePort+id))
}

func StartWork() {

	var err error
	var allPendingClientRequestGuid []string

	err = (*node.GetSystemCoordinator()).WaitUntilLeader(node.GetPropertyAsString(property.WordCountRequestRPCFullAddress))
	utility.CheckError(err)

	(*node.GetLogger()).PrintInfoTaskMessage("Initialization", "I'm leader")

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	allPendingClientRequestGuid, err = (*node.GetSystemCoordinator()).GetAllPendingClientRequestGuid()
	utility.CheckError(err)

	for _, request := range allPendingClientRequestGuid {
		go wordcount.ManageClientRequest(request)
	}

	aftmapreduce.StartAcceptingRPCRequest(&wordcount.Request{}, node.GetPropertyAsString(property.WordCountRequestRPCFullAddress))
}
