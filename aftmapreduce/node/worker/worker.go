package worker

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount"
	"fmt"
)

func Initialize(id int, groupId int, internetAddress string, zookeeperAddresses []string) {

	utility.CheckError(node.InitializeWorker(id, zookeeperAddresses))

	node.SetProperty(property.NodeID, id)
	node.SetProperty(property.NodeGroupID, groupId)
	node.SetProperty(property.NodeType, "Worker")
	node.SetProperty(property.InternetAddress, internetAddress)

	node.SetProperty(property.WordCountMapRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountMapTaskRPCBasePort+id))
	node.SetProperty(property.WordCountReduceRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountReduceTaskRPCBasePort+id))
	node.SetProperty(property.WordCountReceiveRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountReceiveRPCBasePort+id))
	node.SetProperty(property.WordCountSendRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountSendRPCBasePort+id))
	node.SetProperty(property.WordCountRetrieveRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountRetrieverRPCBasePort+id))

}

func StartWork() {

	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Reduce{}, node.GetPropertyAsString(property.WordCountReduceRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Receive{}, node.GetPropertyAsString(property.WordCountReceiveRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Send{}, node.GetPropertyAsString(property.WordCountSendRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Retrieve{}, node.GetPropertyAsString(property.WordCountRetrieveRPCFullAddress))

	err := (*node.GetSystemCoordinator()).RegisterNewWorkerProcess(node.GetPropertyAsInteger(property.NodeID), node.GetPropertyAsInteger(property.NodeGroupID), node.GetPropertyAsString(property.InternetAddress))
	utility.CheckError(err)

	aftmapreduce.StartAcceptingRPCRequest(&wordcount.Map{}, node.GetPropertyAsString(property.WordCountMapRPCFullAddress))
}
