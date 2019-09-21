package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"errors"
)

const (
	AcceptanceJobRequestType        = uint8(0)
	UploadPreSignedURLRequestType   = uint8(1)
	DownloadPreSignedURLRequestType = uint8(2)
)

type Request struct {
}

type RequestInput struct {
	Type             uint8
	SourceFileDigest string
}

type RequestOutput struct {
	Url string
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	var err error
	var isRequestAlreadyAccepted bool

	switch input.Type {
	case UploadPreSignedURLRequestType:
		output.Url, err = (*process.GetStorageKeyValueRegister()).RetrieveURLForPutOperation(input.SourceFileDigest)
	case DownloadPreSignedURLRequestType:
		output.Url, err = (*process.GetStorageKeyValueRegister()).RetrieveURLForGetOperation(input.SourceFileDigest)
	case AcceptanceJobRequestType:

		isRequestAlreadyAccepted, err = (*process.GetSystemCoordinator()).ClientRequestExist(input.SourceFileDigest)
		if err != nil || isRequestAlreadyAccepted {
			break
		} else {
			if err = (*process.GetSystemCoordinator()).RegisterClientRequest(input.SourceFileDigest, acceptancePhaseComplete); err == nil {
				go JobStart(input.SourceFileDigest)
			}
		}
	default:
		return errors.New("request type not recognized")
	}

	return err
}
