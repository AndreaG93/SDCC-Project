package services

type WordCount struct {
}

type WordCountInput struct {
	InputFileName      string
	InputFileDirectory string
}

type WordCountOutput struct {
	OutputFileDigest string
}

func (x *WordCount) Execute(input WordCountInput, output *WordCountOutput) error {
	return nil
}
