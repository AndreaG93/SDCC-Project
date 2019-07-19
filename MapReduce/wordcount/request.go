package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/amazon"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/wordcountfile"
	"SDCC-Project-WorkerNode/utility"
	"net/rpc"
)

type Request struct {
}

type RequestInput struct {
	inputFile string
}

type RequestOutput struct {
	OutputFileDigest string
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	var inputFile *wordcountfile.WordCountFile
	var err error

	if inputFile, err = wordcountfile.New(input.InputFileDigest); err != nil {
		return nil
	}

	//if err = inputFile.DownloadFromCloud(); err != nil {
	//	return nil
	//}

	if err = inputFile.Split(); err != nil {
		return nil
	}

	return nil
}

func SendRequest(inputFilename string, actualPrimaryAddressNode string) {

	var inputFileDigest string
	var err error
	var rpcClient *rpc.Client

	inputFileDigest, err = utility.GenerateDigestOfFileUsingSHA512(inputFilename)
	utility.CheckError(err)

	(*amazon.New()).Upload(inputFilename, inputFileDigest)
	utility.CheckError(err)

	rpcClient, err = rpc.Dial(system.DefaultNetwork, actualPrimaryAddressNode)
	utility.CheckError(err)

	requestOutput := new(RequestOutput)

	err = rpcClient.Call("Request.Execute", &RequestInput{inputFileDigest}, requestOutput)
	utility.CheckError(err)
}
