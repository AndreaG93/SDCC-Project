package wordtokenhashtable

import (
	"testing"
)

func TestWordTokenHashTable_InsertWord(t *testing.T) {

	data := New(10)

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

func TestWordTokenHashTable_Serialize(t *testing.T) {

	var data *WordTokenHashTable
	var dataSerialized []byte
	var dataDeserialized *WordTokenHashTable
	var err error

	data = New(10)

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

	if dataSerialized, err = (*data).Serialize(); err != nil {
		panic(err)
	}

	if dataDeserialized, err = Deserialize(dataSerialized); err != nil {
		panic(err)
	}

	data.Print()
	dataDeserialized.Print()
}
