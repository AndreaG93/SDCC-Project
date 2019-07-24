package utility

import (
	"bytes"
	"encoding/gob"
)

func Encode(data interface{}) ([]byte, error) {

	var err error
	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	if err = enc.Encode(data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Decode(inputData []byte, outputType interface{}) error {

	var err error

	decoder := gob.NewDecoder(bytes.NewReader(inputData))

	if err = decoder.Decode(outputType); err != nil {
		return err
	}

	return nil
}

func MatrixToArray(data [][]byte) []byte {

	dataStructureLength := len(data)
	output := make([]byte, 1)

	output[0] = byte(dataStructureLength)

	for _, dataUnit := range data {

		length := len(dataUnit)
		output = append(output, byte(length))
		output = append(output, dataUnit[:]...)
	}

	return output
}

func ArrayToMatrix(data []byte) [][]byte {

	dataStructureLength := int(data[0])
	output := make([][]byte, dataStructureLength)

	outputIndex := 0
	for index := 1; index < len(data); {

		length := int(data[index])

		subData := data[(index + 1) : (index+1)+length]
		output[outputIndex] = subData
		outputIndex++
		index = index + length + 1
	}

	return output
}
