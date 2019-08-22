package primary

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/wordcount"
	"encoding/gob"
	"fmt"
)

func Initialize(id int, internetAddress string, zookeeperAddresses []string) {

	node.Initialize(zookeeperAddresses)

	node.SetProperty(property.NodeID, id)
	node.SetProperty(property.NodeType, "Primary")
	node.SetProperty(property.InternetAddress, internetAddress)
	node.SetProperty(property.WordCountRequestRPCFullAddress, fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.WordCountRequestRPCBasePort+id))
}

func StartWork() {

	isLeader := make(chan bool)

	go node.GetZookeeperClient().RunAsLeaderCandidate(isLeader, node.GetPropertyAsString(property.WordCountRequestRPCFullAddress))

	<-isLeader

	node.GetLogger().PrintInfoTaskMessage("Initialization", "I'm leader")

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})

	go node.GetZookeeperClient().KeepConnectionAlive()
	wordcount.InitNeededZNodePathsToManageClientRequests()

	for _, request := range wordcount.GetPendingClientsRequests() {
		go wordcount.ManageRequest(request)
	}

	aftmapreduce.StartAcceptingRPCRequest(&wordcount.Request{}, node.GetPropertyAsString(property.WordCountRequestRPCFullAddress))
}
