package data_structures

import (
	"core/utility"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Inserting(t *testing.T) {

	data := BuildWordTokenHashTable(10)

	if err := (*data).InsertWord("Andrea"); err != nil {
		panic(err)
	}
	if err := (*data).InsertWord("Andrea"); err != nil {
		panic(err)
	}
	if err := (*data).InsertWord("Graziani"); err != nil {
		panic(err)
	}
	if err := (*data).InsertWord("Andrea"); err != nil {
		panic(err)
	}

	(*data).Print()
}

func Test_DigestGenerating(t *testing.T) {

	var data *WordTokenHashTable
	var deserializedData *WordTokenHashTable
	var serializedData WordTokenHashTableSerialized
	var digestBeforeSerialization string
	var digestAfterSerialization string
	var err error

	data = BuildWordTokenHashTable(10)

	if err := (*data).InsertWord("Andrea"); err != nil {
		panic(err)
	}
	if err := (*data).InsertWord("Andrea"); err != nil {
		panic(err)
	}
	if err := (*data).InsertWord("Graziani"); err != nil {
		panic(err)
	}
	if err := (*data).InsertWord("Andrea"); err != nil {
		panic(err)
	}

	serializedData = (*data).Serialize()
	digestBeforeSerialization, _ = utility.SHA512(serializedData)

	if err = utility.WriteToLocalDisk(digestBeforeSerialization, serializedData); err != nil {
		panic(err)
	}

	if serializedData, err = ReadWordTokenHashTableSerializedFromLocalDisk(digestBeforeSerialization); err != nil {
		panic(err)
	}

	digestAfterSerialization, _ = utility.SHA512(serializedData)
	fmt.Println(digestBeforeSerialization)
	fmt.Println(digestAfterSerialization)

	if strings.Compare(digestBeforeSerialization, digestAfterSerialization) != 0 {
		log.Fatal("Digest NOT correct!")
	}

	if deserializedData, err = serializedData.Deserialize(); err != nil {
		panic(err)
	}

	deserializedData.Print()
}
