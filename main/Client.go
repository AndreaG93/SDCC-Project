package main

import (
	"SDCC-Project-WorkerNode/cloud/amazon/s3"
	"SDCC-Project-WorkerNode/core"
	"SDCC-Project-WorkerNode/services"
	"SDCC-Project-WorkerNode/utility"
	"net/rpc"
)

const (
	BUCKET_NAME = "graziani-filestorage"
)

func main() {

	var inputFileDigest string
	var err error
	var rpcClient *rpc.Client
	var s3Client *s3.Client
	var wordCountInput *services.WordCountInput
	var wordCountOutput *services.WordCountOutput

	wordCountInput = new(services.WordCountInput)
	wordCountOutput = new(services.WordCountOutput)

	if inputFileDigest, err = utility.GenerateDigestOfFileUsingSHA512("./src/test-input-data/input.txt"); err != nil {
		panic(err)
	}

	s3Client = s3.New(core.DefaultAmazonAWSRegion, BUCKET_NAME)

	err = s3Client.Upload("./src/test-input-data/input.txt", inputFileDigest)
	utility.CheckError(err)

	rpcClient, err = rpc.Dial(core.DefaultNetwork, "localhost:5000")
	utility.CheckError(err)

	wordCountInput.InputFileName = inputFileDigest
	wordCountInput.InputFileDirectory = BUCKET_NAME

	err = rpcClient.Call("WordCount.Execute", wordCountInput, wordCountOutput)
	utility.CheckError(err)

	err = s3Client.Download("./output", inputFileDigest)
	utility.CheckError(err)

}
