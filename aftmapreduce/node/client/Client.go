package client

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"encoding/gob"
	"fmt"
	"net/rpc"
)

func StartWork(sourceFilePath string, zookeeperAddresses []string) {

	initialize(zookeeperAddresses)

	currentLeaderInternetAddress, err := node.GetZookeeperClient().GetCurrentLeaderRequestRPCInternetAddress()
	utility.CheckError(err)

	sourceFileDigest, err := utility.GenerateDigestOfFileUsingSHA512(sourceFilePath)
	utility.CheckError(err)

	node.GetAmazonS3Client().Upload(sourceFilePath, sourceFileDigest)

	rawDataOutput := sendRequestToCurrentLeader(sourceFileDigest, currentLeaderInternetAddress)
	printResult(rawDataOutput)
}

func initialize(zookeeperAddresses []string) {
	node.Initialize(zookeeperAddresses)

	gob.Register(wordcount.Request{})
}

func sendRequestToCurrentLeader(sourceFileDigest string, currentLeaderInternetAddress string) []byte {

	input := new(wordcount.RequestInput)
	output := new(wordcount.RequestOutput)

	(*input).SourceFileDigest = sourceFileDigest

	client, err := rpc.Dial("tcp", currentLeaderInternetAddress)
	utility.CheckError(err)

	err = client.Call("Request.Execute", &input, &output)
	utility.CheckError(err)

	finalOutputPath := fmt.Sprintf("%s/%s", wordcount.CompleteRequestsZNodePath, sourceFileDigest)

	rawData, watcher := node.GetZookeeperClient().GetZNodeData(finalOutputPath)
	if rawData == nil {
		<-watcher
		rawData, _ = node.GetZookeeperClient().GetZNodeData(finalOutputPath)
	}

	return rawData
}

func printResult(rawData []byte) {
	result, err := WordTokenList.Deserialize(rawData)
	utility.CheckError(err)

	result.Print()
}
