package client

import (
	"SDCC-Project/aftmapreduce/storage"
	"SDCC-Project/aftmapreduce/system/zookeeper"
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

	zookeeperClient, err = zookeeper.New(zookeeperAddresses)
	utility.CheckError(err)

	data, err := ioutil.ReadFile(sourceFilePath)
	utility.CheckError(err)

	sourceFileDigest := utility.GenerateDigestUsingSHA512(data)
	utility.CheckError(err)

	for {

		currentLeaderInternetAddress, err = zookeeperClient.GetCurrentLeaderPublicInternetAddress()
		fmt.Println("Current Leader IP: " + currentLeaderInternetAddress)
		if err != nil {
			continue
		}

		preSignedURL, err = sendRequest(sourceFileDigest, currentLeaderInternetAddress, wordcount.UploadPreSignedURLRequestType)
		if err != nil {
			continue
		} else {
			break
		}
	}

	for {
		err = storage.Upload(preSignedURL, data)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			break
		}
	}

	for {

		currentLeaderInternetAddress, err = zookeeperClient.GetCurrentLeaderPublicInternetAddress()
		if err != nil {
			continue
		}

		_, err := sendRequest(sourceFileDigest, currentLeaderInternetAddress, wordcount.AcceptanceJobRequestType)
		if err != nil {
			continue
		} else {
			break
		}
	}

	fmt.Println("Waiting for job completion...")
	for {
		outputDigest, err = zookeeperClient.WaitForClientRequestCompletion(sourceFileDigest)
		if err == nil {
			break
		}
	}
	fmt.Println("Job completion signal received...")

	for {

		currentLeaderInternetAddress, err = zookeeperClient.GetCurrentLeaderPublicInternetAddress()
		if err != nil {
			continue
		}

		preSignedURL, err = sendRequest(outputDigest, currentLeaderInternetAddress, wordcount.DownloadPreSignedURLRequestType)
		if err != nil {
			continue
		} else {
			break
		}
	}

	for {
		outputBytes, err := storage.Download(preSignedURL)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			printResult(outputBytes)
			return
		}
	}

}

func sendRequest(sourceFileDigest string, currentLeaderInternetAddress string, requestType uint8) (string, error) {

	input := new(wordcount.RequestInput)
	output := new(wordcount.RequestOutput)

	(*input).SourceFileDigest = sourceFileDigest
	(*input).Type = requestType

	client, err := rpc.Dial("tcp", currentLeaderInternetAddress)
	if err != nil {
		return "", err
	}

	err = client.Call("Request.Execute", &input, &output)
	if err != nil {
		return "", err
	}

	return output.Url, nil
}

func printResult(rawData []byte) {
	if result, err := WordTokenList.Deserialize(rawData); err != nil {
		panic(err)
	} else {
		result.Print()
	}
}
