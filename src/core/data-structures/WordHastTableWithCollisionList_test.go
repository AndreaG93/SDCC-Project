package data_structures

import (
	"container/list"
	"core/utility"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Inserting(t *testing.T) {

	const (
		input1 = "Andrea"
		input2 = "Graziani"
	)

	data := BuildWordHashTableWithCollisionList(10)

	if err := data.Insert(input1); err != nil {
		log.Fatal("Insert Operation Failed!")
	}
	if err := data.Insert(input2); err != nil {
		log.Fatal("Insert Operation Failed!")
	}
	if err := data.Insert(input1); err != nil {
		log.Fatal("Insert Operation Failed!")
	}
	if err := data.Insert(input1); err != nil {
		log.Fatal("Insert Operation Failed!")
	}

	var index1 uint
	var index2 uint
	var output1 *WordToken
	var output2 *WordToken
	var currentList1 *list.List
	var currentList2 *list.List
	var err error

	if index1, err = utility.GenerateArrayIndexFromString(input1, uint(10)); err != nil {
		log.Fatal("Fail while generating index!")
	}

	if index2, err = utility.GenerateArrayIndexFromString(input2, uint(10)); err != nil {
		log.Fatal("Fail while generating index!")
	}

	data.Print()

	currentList1 = data.getCollisionListAt(index1)
	currentList2 = data.getCollisionListAt(index2)

	output1 = currentList1.Front().Value.(*WordToken)
	output2 = currentList2.Front().Value.(*WordToken)

	if strings.Compare((*output1).Word, input1) != 0 {
		log.Fatal("Output 1: NOT correct!")
	}

	if (*output1).Occurrences != 3 {
		log.Fatal("Output 1: NOT correct!")
	}

	if strings.Compare(output2.Word, input2) != 0 {
		log.Fatal("Output 2: NOT correct!")
	}

	if output2.Occurrences != 1 {
		log.Fatal("Output 2: NOT correct!")
	}

	fmt.Println(data.numberOfOccurrences)
	if data.numberOfOccurrences != 2 {
		log.Fatal("Number of stored occurrences NOT correct!")
	}
}

func Test_DigestGenerating(t *testing.T) {

	const (
		input1         = "Andrea"
		input2         = "Graziani"
		expectedOutput = "7cf06db9ae13658cd35c9e36433b09ef80682de9af846cdde9645c7fc477341391418c762ef172bb42d0f3419bb3fdbbca1d9687c5d17d4549cc6c99f79e305d"
	)

	data := BuildWordHashTableWithCollisionList(10)

	if err := data.Insert(input1); err != nil {
		log.Fatal("Insert Operation Failed!")
	}
	if err := data.Insert(input2); err != nil {
		log.Fatal("Insert Operation Failed!")
	}
	if err := data.Insert(input1); err != nil {
		log.Fatal("Insert Operation Failed!")
	}
	if err := data.Insert(input1); err != nil {
		log.Fatal("Insert Operation Failed!")
	}

	var digest string

	digest, _ = data.GetDigest()

	fmt.Println(digest)

	if strings.Compare(digest, expectedOutput) != 0 {
		log.Fatal("Digest NOT correct!")
	}
}
