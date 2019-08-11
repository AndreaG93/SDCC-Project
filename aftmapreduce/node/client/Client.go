package client

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"encoding/gob"
	"net/rpc"
)

func StartWork(sourceFilePath string, zookeeperAddresses []string) {

	initialize(zookeeperAddresses)

	currentLeaderInternetAddress, err := node.GetZookeeperClient().GetCurrentLeaderInternetAddress()
	utility.CheckError(err)

	sourceFileDigest, err := utility.GenerateDigestOfFileUsingSHA512(sourceFilePath)
	utility.CheckError(err)

	node.GetAmazonS3Client().Upload(sourceFilePath, sourceFileDigest)

	sendRequestToCurrentLeader(sourceFileDigest, currentLeaderInternetAddress)
}

func initialize(zookeeperAddresses []string) {
	node.Initialize(zookeeperAddresses)

	gob.Register(wordcount.Request{})
}

func sendRequestToCurrentLeader(sourceFileDigest string, currentLeaderInternetAddress string) {

	input := new(wordcount.RequestInput)
	output := new(wordcount.RequestOutput)

	(*input).SourceFileDigest = sourceFileDigest

	client, err := rpc.Dial("tcp", currentLeaderInternetAddress)
	utility.CheckError(err)

	err = client.Call("Request.Execute", &input, &output)
	utility.CheckError(err)
}

func printResult(rawData []byte) {
	result, err := WordTokenList.Deserialize(rawData)
	utility.CheckError(err)

	result.Print()
}

/*
func kkStartWork(filename string, zookeeperAddresses []string) {

	var rawData []byte
	var watcher <-chan zk.Event

	digest, err := utility.GenerateDigestOfFileUsingSHA512(filename)
	if err != nil {
		panic(err)
	}

	node.Initialize(0, "Client", zookeeperAddresses)

	path := fmt.Sprintf("%s/%s", aftmapreduce.CompleteRequestsZNodePath, digest)

	if !node.GetZookeeperClient().CheckZNodeExistence(path) {

		S3Client := amazons3.New()
		S3Client.Upload(filename, digest)

		internetAddress, err := node.GetZookeeperClient().GetCurrentLeaderInternetAddress()
		if err != nil {
			panic(err)
		}

		sendRequest(digest, internetAddress)

		_, watcher = node.GetZookeeperClient().GetZNodeData(path)
		<-watcher
	} else {

		internetAddress, err := node.GetZookeeperClient().GetCurrentLeaderInternetAddress()
		if err != nil {
			panic(err)
		}

		rawData, _ = node.GetZookeeperClient().GetZNodeData(path)
		if rawData == nil {
			sendRequest(digest, internetAddress)
			_, watcher = node.GetZookeeperClient().GetZNodeData(path)
			<-watcher
		}
	}

	rawData, _ = node.GetZookeeperClient().GetZNodeData(path)
	printResult(rawData)
}*/
