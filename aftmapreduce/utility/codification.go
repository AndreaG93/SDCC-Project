package utility

import (
	"bytes"
	"encoding/gob"
)

func Encode(data interface{}) []byte {

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	CheckError(encoder.Encode(data))

	return buffer.Bytes()
}

func Decode(inputData []byte, outputType interface{}) {

	decoder := gob.NewDecoder(bytes.NewReader(inputData))
	CheckError(decoder.Decode(outputType))
}

func Decoding(rawData []byte, outputType interface{}) error {

	decoder := gob.NewDecoder(bytes.NewReader(rawData))
	return decoder.Decode(outputType)
}

func Encoding(data interface{}) ([]byte, error) {

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	} else {
		return buffer.Bytes(), nil
	}
}
