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

	output, _ := GenerateDigestOfDataUsingSHA512(input)

	if strings.Compare(output, expectedDigest) != 0 {
		log.Fatal("Output NOT correct!")
	}
}

func TestGenerateDigestOfDataUsingSHA5122(t *testing.T) {

	const (
		input          = "Andrea Graziani"
		expectedDigest = "33d70a373d75aea143cdeff350c48f7c51cc7134247ad15758edb56ddd09c7bdcb531e01a7f7e006dae7c0f2be765b558e5583c11f86f6084fb3341937fc7117"
	)

	output := test(input)

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

func BenchmarkGenerateDigestOfDataUsingSHA5122(b *testing.B) {

	for n := 0; n < b.N; n++ {
		TestGenerateDigestOfDataUsingSHA5122(nil)
	}
}
