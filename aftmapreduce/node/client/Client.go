package client

import (
	"SDCC-Project/aftmapreduce/cloud/amazons3"
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

	var err error
	var currentLeaderInternetAddress string
	var preSignedURL string
	var outputDigest string

	zookeeperClient = zookeeper.New(zookeeperAddresses)

	data, err := ioutil.ReadFile(sourceFilePath)
	utility.CheckError(err)

	sourceFileDigest := utility.GenerateDigestUsingSHA512(data)
	utility.CheckError(err)

	for {

		currentLeaderInternetAddress, err = zookeeperClient.GetCurrentLeaderRequestRPCInternetAddress()
		if err != nil {
			continue
		}

		preSignedURL, err = sendRequestForPreSignedURL(sourceFileDigest, currentLeaderInternetAddress, true, false)
		if err != nil {
			continue
		} else {
			break
		}
	}

	for {
		err = amazons3.UploadWithPreSignedURL(data, preSignedURL)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			break
		}
	}

	for {

		currentLeaderInternetAddress, err = zookeeperClient.GetCurrentLeaderRequestRPCInternetAddress()
		if err != nil {
			continue
		}

		err := sendRequestForJob(sourceFileDigest, currentLeaderInternetAddress)
		if err != nil {
			continue
		} else {
			break
		}
	}

	finalOutputPath := fmt.Sprintf("%s/%s", wordcount.CompleteRequestsZNodePath, sourceFileDigest)

	for {

		rawData, watcher := zookeeperClient.GetZNodeData(finalOutputPath)
		if rawData == nil {
			<-watcher
			rawData, _ = zookeeperClient.GetZNodeData(finalOutputPath)
		} else {
			outputDigest = string(rawData)
			break
		}
	}

	for {

		currentLeaderInternetAddress, err = zookeeperClient.GetCurrentLeaderRequestRPCInternetAddress()
		if err != nil {
			continue
		}

		preSignedURL, err = sendRequestForPreSignedURL(outputDigest, currentLeaderInternetAddress, false, true)
		if err != nil {
			continue
		} else {
			break
		}
	}

	for {
		outputBytes, err := amazons3.DownloadPreSignedURL(preSignedURL)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			printResult(outputBytes)
			return
		}
	}

}

func sendRequestForPreSignedURL(sourceFileDigest string, currentLeaderInternetAddress string, upload bool, download bool) (string, error) {

	input := new(wordcount.RequestInput)
	output := new(wordcount.RequestOutput)

	(*input).SourceFileDigest = sourceFileDigest
	(*input).RequestPreSignedURLForUpload = upload
	(*input).RequestPreSignedURLForDownload = download

	client, err := rpc.Dial("tcp", currentLeaderInternetAddress)
	if err != nil {
		return "", err
	}

	err = client.Call("Request.Execute", &input, &output)
	if err != nil {
		return "", err
	}

	return output.PreSignedURL, nil
}

func sendRequestForJob(sourceFileDigest string, currentLeaderInternetAddress string) error {

	input := new(wordcount.RequestInput)
	output := new(wordcount.RequestOutput)

	(*input).SourceFileDigest = sourceFileDigest
	(*input).RequestPreSignedURLForDownload = false
	(*input).RequestPreSignedURLForUpload = false

	client, err := rpc.Dial("tcp", currentLeaderInternetAddress)
	if err != nil {
		return err
	}

	err = client.Call("Request.Execute", &input, &output)
	if err != nil {
		return err
	}

	return nil
}

func printResult(rawData []byte) {
	result := WordTokenList.Deserialize(rawData)
	result.Print()
}
