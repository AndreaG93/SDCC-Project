package utility

import (
	"log"
	"strings"
	"testing"
)

func TestGenerateDigestOfDataUsingSHA512(t *testing.T) {

	const (
		input          = "Andrea Graziani"
		expectedDigest = "61c6df4cf5e8ca09a7c03f10587190c0d8b0609c3c145d37390f740939b92c70ef8f517f4f99457c11007d943cead69349a8aad413b662fec6efab6932a2aeb5"
	)

	output, _ := GenerateDigestUsingSHA512([]byte(input))

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

func TestGenerateFileDigestUsingSHA512(t *testing.T) {

	const (
		expectedDigest = "b93b18bd188e4da4b1ad7bc0807a9cd40ef75136469c11e848ab0336001145a80df5532f6467a9ee0a42e90fe40e217ccf152ddfaebdc9623ff306a9e4270816"
	)

	var digest string
	var err error

	if digest, err = GenerateDigestOfFileUsingSHA512("../../test-input-data/input.txt"); err != nil {
		panic(err)
	}

	if strings.Compare(expectedDigest, digest) != 0 {
		log.Fatal("Output NOT correct!")
	}
}

func BenchmarkGenerateFileDigestUsingSHA512(b *testing.B) {

	for n := 0; n < b.N; n++ {
		TestGenerateFileDigestUsingSHA512(nil)
	}
}

func BenchmarkGenerateDigestOfDataUsingSHA512(b *testing.B) {

	for n := 0; n < b.N; n++ {
		TestGenerateFileDigestUsingSHA512(nil)
	}
}
