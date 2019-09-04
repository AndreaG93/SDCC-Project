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
