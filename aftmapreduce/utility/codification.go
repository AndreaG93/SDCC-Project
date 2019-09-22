package utility

import (
	"bytes"
	"encoding/gob"
)

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
