package services

type Map struct {
}

type MapInput struct {
	InputString                  string
	OutputWordHashTableArraySize uint
}

type MapOutput struct {
	OutputFileDigest string
}

func (x *Map) Execute(input MapInput, output *MapOutput) error {

	/*
		var outputData *data_structures.WordTokenHashTable
		var outputDataDigest string
		var err error

		outputData = data_structures.BuildWordTokenHashTable(input.OutputWordHashTableArraySize)

		wordScanner := utility.BuildWordScannerFromString(input.InputString)

		for wordScanner.Scan() {

			currentWord := strings.ToLower(wordScanner.Text())
			if err = outputData.InsertWord(currentWord); err != nil {
				return err
			}
		}

		if outputDataDigest, err = outputData.GetDigest(); err != nil {
			return err
		}

		if err = outputData.WriteOnLocalDisk(); err != nil {
			return err
		}

		output.OutputFileDigest = outputDataDigest
	*/
	return nil
}
