package wordcount

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/wordcountfile"
)

type Request struct {
}

type RequestInput struct {
	InputFileName   string
	InputFileDigest string
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

	if err = inputFile.DownloadFromCloud(); err != nil {
		return nil
	}

	if err = inputFile.Split(); err != nil {
		return nil
	}

	return nil
}
