package WordTokenHashTable

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
