package wordtokenlist

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtoken"
	"log"
	"strings"
	"testing"
)

func Test_WordTokenList_Inserting(t *testing.T) {

	var currentWordToken *wordtoken.WordToken

	wordTokenList := New()

	(*wordTokenList).InsertWordToken(wordtoken.New("Andrea", 5))
	(*wordTokenList).InsertWordToken(wordtoken.New("Graziani", 30))
	(*wordTokenList).InsertWordToken(wordtoken.New("Yumi", 26))
	(*wordTokenList).InsertWordToken(wordtoken.New("Yumi", 4))
	(*wordTokenList).InsertWordToken(wordtoken.New("Andrea", 5))

	(*wordTokenList).Print()

	(*wordTokenList).Next()
	currentWordToken = (*wordTokenList).WordToken()

	if strings.Compare((*currentWordToken).Word, "Andrea") != 0 {
		log.Fatal("NOT correct!")
	}
	if (*currentWordToken).Occurrences != 10 {
		log.Fatal("NOT correct!")
	}

	(*wordTokenList).Next()
	currentWordToken = (*wordTokenList).WordToken()

	if strings.Compare((*currentWordToken).Word, "Graziani") != 0 {
		log.Fatal("NOT correct!")
	}
	if (*currentWordToken).Occurrences != 30 {
		log.Fatal("NOT correct!")
	}

	(*wordTokenList).Next()
	currentWordToken = (*wordTokenList).WordToken()

	if strings.Compare((*currentWordToken).Word, "Yumi") != 0 {
		log.Fatal("NOT correct!")
	}
	if (*currentWordToken).Occurrences != 30 {
		log.Fatal("NOT correct!")
	}
}

func TestWordTokenList_Serialize(t *testing.T) {

	var data *WordTokenList
	var dataSerialized []byte
	var dataDeserialized *WordTokenList
	var err error

	data = New()

	(*data).InsertWordToken(wordtoken.New("Andrea", 5))
	(*data).InsertWordToken(wordtoken.New("Graziani", 5))
	(*data).InsertWordToken(wordtoken.New("Diana", 5))

	if dataSerialized, err = (*data).Serialize(); err != nil {
		panic(err)
	}

	if dataDeserialized, err = Deserialize(dataSerialized); err != nil {
		panic(err)
	}

	data.Print()
	dataDeserialized.Print()
}

func Test_WordTokenList_WordTokenListMerging(t *testing.T) {

	var currentWordToken *wordtoken.WordToken

	wordTokenList1 := New()
	wordTokenList2 := New()

	(*wordTokenList1).InsertWordToken(wordtoken.New("Andrea", 5))
	(*wordTokenList1).InsertWordToken(wordtoken.New("Graziani", 30))
	(*wordTokenList1).InsertWordToken(wordtoken.New("Yumi", 26))
	(*wordTokenList2).InsertWordToken(wordtoken.New("Andrea", 5))
	(*wordTokenList2).InsertWordToken(wordtoken.New("Graziani", 30))
	(*wordTokenList2).InsertWordToken(wordtoken.New("Yumi", 26))

	(*wordTokenList1).Merge(wordTokenList2)

	(*wordTokenList1).Next()
	currentWordToken = (*wordTokenList1).WordToken()

	if strings.Compare((*currentWordToken).Word, "Andrea") != 0 {
		log.Fatal("NOT correct!")
	}
	if (*currentWordToken).Occurrences != 10 {
		log.Fatal("NOT correct!")
	}

	(*wordTokenList1).Next()
	currentWordToken = (*wordTokenList1).WordToken()

	if strings.Compare((*currentWordToken).Word, "Graziani") != 0 {
		log.Fatal("NOT correct!")
	}
	if (*currentWordToken).Occurrences != 60 {
		log.Fatal("NOT correct!")
	}

	(*wordTokenList1).Next()
	currentWordToken = (*wordTokenList1).WordToken()

	if strings.Compare((*currentWordToken).Word, "Yumi") != 0 {
		log.Fatal("NOT correct!")
	}
	if (*currentWordToken).Occurrences != 52 {
		log.Fatal("NOT correct!")
	}
}
