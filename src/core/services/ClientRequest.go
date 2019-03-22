package services

type ClientRequest struct {
}

type ClientRequestInput struct {
	inputFileName string
}

type ClientRequestOutput struct {
	OutputFileDigest string
}

func (x *ClientRequest) Execute(input MapInput, output *MapOutput) error {
	return nil
}
