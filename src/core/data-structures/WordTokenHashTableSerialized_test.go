package data_structures

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_SerializationDeserialization(t *testing.T) {

	var err error
	var wordTokenHashTable *WordTokenHashTable
	var digest string
	var digestOfDataRetrievedFromDisk string

	const (
		input1 = "Andrea"
		input2 = "Graziani"
	)

	wordTokenHashTable = BuildWordTokenHashTable(10)

	if err := wordTokenHashTable.InsertWord(input1); err != nil {
		panic(err)
	}
	if err := wordTokenHashTable.InsertWord(input2); err != nil {
		panic(err)
	}
	if err := wordTokenHashTable.InsertWord(input1); err != nil {
		panic(err)
	}
	if err := wordTokenHashTable.InsertWord(input1); err != nil {
		panic(err)
	}

	if digest, err = wordTokenHashTable.GetDigest(); err != nil {
		panic(err)
	}

	fmt.Printf("Digest of Orginal Data:\n%s\n", digest)

	if err = wordTokenHashTable.WriteOnLocalDisk(); err != nil {
		panic(err)
	}

	if wordTokenHashTable, err = ReadTokenWordHashTableFromLocalDisk(digest); err != nil {
		panic(err)
	}

	if digestOfDataRetrievedFromDisk, err = wordTokenHashTable.GetDigest(); err != nil {
		panic(err)
	}

	fmt.Printf("Digest of wordTokenList read from disk:\n%s\n", digestOfDataRetrievedFromDisk)

	if strings.Compare(digest, digestOfDataRetrievedFromDisk) != 0 {
		log.Fatal("FAILED: Different digest")
	}

	wordTokenHashTable.Print()
}
