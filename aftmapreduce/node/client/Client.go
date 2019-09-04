package client

import (
	"SDCC-Project/aftmapreduce/cloud/zookeeper"
	"SDCC-Project/aftmapreduce/utility"
	"SDCC-Project/aftmapreduce/wordcount"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"fmt"
	"io/ioutil"
	"net/rpc"
)

var zookeeperClient *zookeeper.Client

func StartWork(sourceFilePath string, zookeeperAddresses []string) {

	zookeeperClient = zookeeper.New(zookeeperAddresses)

	currentLeaderInternetAddress, err := zookeeperClient.GetCurrentLeaderRequestRPCInternetAddress()
	utility.CheckError(err)

	data, err := ioutil.ReadFile(sourceFilePath)
	utility.CheckError(err)

	sourceFileDigest := utility.GenerateDigestUsingSHA512(data)
	utility.CheckError(err)

	rawDataOutput := sendRequestToCurrentLeader(sourceFileDigest, string(data), currentLeaderInternetAddress)
	printResult(rawDataOutput)
}

func sendRequestToCurrentLeader(sourceFileDigest string, fileContent string, currentLeaderInternetAddress string) []byte {

	input := new(wordcount.RequestInput)
	output := new(wordcount.RequestOutput)

	(*input).SourceFileDigest = sourceFileDigest
	(*input).FileContent = fileContent

	client, err := rpc.Dial("tcp", currentLeaderInternetAddress)
	utility.CheckError(err)

	err = client.Call("Request.Execute", &input, &output)
	utility.CheckError(err)

	finalOutputPath := fmt.Sprintf("%s/%s", wordcount.CompleteRequestsZNodePath, sourceFileDigest)

	rawData, watcher := zookeeperClient.GetZNodeData(finalOutputPath)
	if rawData == nil {
		<-watcher
		rawData, _ = zookeeperClient.GetZNodeData(finalOutputPath)
	}

	return rawData
}

func printResult(rawData []byte) {
	result := WordTokenList.Deserialize(rawData)
	result.Print()
}
