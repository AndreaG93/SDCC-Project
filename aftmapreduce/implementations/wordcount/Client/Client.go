package Client

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/implementations/wordcount"
	"SDCC-Project/aftmapreduce/implementations/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/cloud/amazon"
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/utility"
	"encoding/gob"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"net/rpc"
)

func sendRequest(digestFile string, internetAddress string) {

	input := new(aftmapreduce.EntryPointInput)
	output := new(aftmapreduce.EntryPointOutput)

	inputData := new(wordcount.Input)
	(*inputData).FileDigest = digestFile
	(*inputData).MapCardinality = 5

	input.Data = inputData

	gob.Register(wordcount.Input{})

	currentLeaderPublicInternetAddress := fmt.Sprintf("%s:%d", internetAddress, aftmapreduce.MapReduceRequestRPCBasePort+1)

	client, err := rpc.Dial("tcp", currentLeaderPublicInternetAddress)
	utility.CheckError(err)

	err = client.Call("EntryPoint.Execute", &input, &output)
	utility.CheckError(err)
}

func printResult(rawData []byte) {
	result, err := WordTokenList.Deserialize(rawData)
	utility.CheckError(err)

	result.Print()
}

func StartWork(filename string, internetAddress string) {

	var rawData []byte
	var watcher <-chan zk.Event

	digest, err := utility.GenerateDigestOfFileUsingSHA512(filename)
	if err != nil {
		panic(err)
	}

	zookeeperClient := zookeeper.New([]string{"localhost:2181"})
	path := fmt.Sprintf("%s/%s", aftmapreduce.CompleteRequestsZNodePath, digest)

	if !zookeeperClient.CheckZNodeExistence(path) {

		S3Client := amazon.New()
		S3Client.Upload(filename, digest)

		sendRequest(digest, internetAddress)
		_, watcher = zookeeperClient.GetZNodeData(path)
		<-watcher
	} else {
		rawData, _ = zookeeperClient.GetZNodeData(path)
		if rawData == nil {
			sendRequest(digest, internetAddress)
			_, watcher = zookeeperClient.GetZNodeData(path)
			<-watcher
		}
	}

	rawData, _ = zookeeperClient.GetZNodeData(path)
	printResult(rawData)
}
