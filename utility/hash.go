package utility

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash/fnv"
	"os/exec"
	"sync"
)

// FNV-1a is a not cryptographic hash function:
// 1) Fast to compute and designed for fast hash table.
// 2) Slightly better avalanche characteristics than FNV-1 hash function.
var FNV1AHashAlgorithm = fnv.New32a()

// Secure Hash Algorithm...
var SHA512cryptoHashAlgorithm = sha512.New()

var mutex = &sync.Mutex{}

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

func GenerateDigestUsingSHA512(data []byte) string {

	mutex.Lock()
	_, err := SHA512cryptoHashAlgorithm.Write(data)
	CheckError(err)

	output := hex.EncodeToString(SHA512cryptoHashAlgorithm.Sum(nil))
	SHA512cryptoHashAlgorithm.Reset()
	mutex.Unlock()

	return output
}

func GenerateDigestOfFileUsingSHA512(filename string) (string, error) {

	var command *exec.Cmd
	var commandOutput []byte
	var commandError error

	command = exec.Command("sha512sum", filename)

	if commandOutput, commandError = command.Output(); commandError != nil {
		return "", commandError
	}

	return string(commandOutput[:128]), nil
}
