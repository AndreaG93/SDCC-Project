package utility

import (
	"log"
	"strings"
	"testing"
)

func TestGenerateDigestOfDataUsingSHA512(t *testing.T) {

	const (
		input          = "Andrea Graziani"
		expectedDigest = "33d70a373d75aea143cdeff350c48f7c51cc7134247ad15758edb56ddd09c7bdcb531e01a7f7e006dae7c0f2be765b558e5583c11f86f6084fb3341937fc7117"
	)

	output := GenerateDigestUsingSHA512([]byte(input))

	if strings.Compare(output, expectedDigest) != 0 {
		log.Fatal("Output NOT correct!")
	}
}

func TestGenerateArrayIndexFromString(t *testing.T) {

	const (
		arrayIndex      = 10
		input1          = "Andrea"
		expectedOutput1 = 2
	)

	output1, _ := GenerateArrayIndexFromString(input1, arrayIndex)

	if output1 != expectedOutput1 {
		log.Fatal("Output NOT correct!")
	}
}
