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
