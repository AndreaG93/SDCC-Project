package utility

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"hash/fnv"
)

// FNV-1a is a not cryptographic hash function:
// 1) Fast to compute and designed for fast hash table.
// 2) Slightly better avalanche characteristics than FNV-1 hash function.
var FNV1AHashAlgorithm = fnv.New32a()

// Secure Hash Algorithm...
var SHA512cryptoHashAlgorithm = sha512.New()

func GenerateArrayIndexFromString(inputString string, arraySize uint) (uint, error) {

	if inputString == "" {
		return 0, errors.New(InvalidInput)
	}

	if _, err := FNV1AHashAlgorithm.Write([]byte(inputString)); err != nil {
		return 0, err
	}
	defer FNV1AHashAlgorithm.Reset()

	return uint(FNV1AHashAlgorithm.Sum32()) % arraySize, nil
}

func SHA512(data interface{}) (string, error) {

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(data); err != nil {
		return "", err
	}

	if _, err := SHA512cryptoHashAlgorithm.Write(buffer.Bytes()); err != nil {
		return "", err
	}
	defer SHA512cryptoHashAlgorithm.Reset()

	return hex.EncodeToString(SHA512cryptoHashAlgorithm.Sum(nil)), nil
}
