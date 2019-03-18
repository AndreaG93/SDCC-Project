package utility

import (
	"log"
	"strings"
	"testing"
)

func Test_SHA512UsingString(t *testing.T) {

	const (
		input1          = "Andrea Graziani"
		input2          = "Unione Europea"
		expectedOutput1 = "61c6df4cf5e8ca09a7c03f10587190c0d8b0609c3c145d37390f740939b92c70ef8f517f4f99457c11007d943cead69349a8aad413b662fec6efab6932a2aeb5"
		expectedOutput2 = "1195445321fa20d7cc9cd9f96f3dd58e7b72132a93cf9a95e3a819ad92af81215e29000af5a83e343737ee8f1de6c2dda6d1da7550258d03d7d019046370a0d2"
	)

	output1, _ := SHA512(input1)
	output2, _ := SHA512(input2)

	if strings.Compare(output1, expectedOutput1) != 0 {
		log.Fatal("Output 1: NOT correct!")
	}

	if strings.Compare(output2, expectedOutput2) != 0 {
		log.Fatal("Output 2: NOT correct!")
	}
}

func Test_ArrayIndexGeneration(t *testing.T) {

	const (
		arrayIndex      = 10
		input1          = "Andrea"
		expectedOutput1 = 2
	)

	output1, _ := GenerateArrayIndexFromString(input1, arrayIndex)

	if output1 != expectedOutput1 {
		log.Fatal("Output 1: NOT correct!")
	}
}
