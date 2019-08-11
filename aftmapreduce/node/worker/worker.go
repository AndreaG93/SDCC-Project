package worker

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/wordcount"
	"encoding/gob"
	"fmt"
)

func Initialize(id int, groupId int, internetAddress string, zookeeperAddresses []string) {

	node.Initialize(zookeeperAddresses)

	node.SetProperty(property.NodeID, id)
	node.SetProperty(property.NodeGroupID, groupId)
	node.SetProperty(property.NodeType, "Worker")
	node.SetProperty(property.InternetAddress, internetAddress)

	node.SetProperty(property.WordCountMapRPCFullAddress, fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.WordCountMapTaskRPCBasePort+id))
	node.SetProperty(property.WordCountReduceRPCFullAddress, fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.WordCountReduceTaskRPCBasePort+id))
	node.SetProperty(property.WordCountReceiveRPCFullAddress, fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.WordCountReceiveRPCBasePort+id))
	node.SetProperty(property.WordCountSendRPCFullAddress, fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.WordCountSendRPCBasePort+id))

	gob.Register(wordcount.MapInput{})
	gob.Register(wordcount.ReduceInput{})
	gob.Register(wordcount.Send{})
}

func StartWork() {

	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Map{}, node.GetPropertyAsString(property.WordCountMapRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Reduce{}, node.GetPropertyAsString(property.WordCountReduceRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Receive{}, node.GetPropertyAsString(property.WordCountReceiveRPCFullAddress))
	go aftmapreduce.StartAcceptingRPCRequest(&wordcount.Send{}, node.GetPropertyAsString(property.WordCountSendRPCFullAddress))

	node.GetZookeeperClient().RegisterNodeMembership(node.GetPropertyAsInteger(property.NodeID), node.GetPropertyAsInteger(property.NodeGroupID), node.GetPropertyAsString(property.InternetAddress))
	node.GetZookeeperClient().KeepConnectionAlive()
}
