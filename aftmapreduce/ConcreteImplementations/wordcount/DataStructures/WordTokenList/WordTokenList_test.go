package WordTokenList

import (
	"SDCC-Project/aftmapreduce/ConcreteImplementations/wordcount/DataStructures/WordToken"
	"log"
	"strings"
	"testing"
)

func Test_WordTokenList_Inserting(t *testing.T) {

	var currentWordToken *WordToken.WordToken

	wordTokenList := New()

	(*wordTokenList).InsertWordToken(WordToken.New("Andrea", 5))
	(*wordTokenList).InsertWordToken(WordToken.New("Graziani", 30))
	(*wordTokenList).InsertWordToken(WordToken.New("Yumi", 26))
	(*wordTokenList).InsertWordToken(WordToken.New("Yumi", 4))
	(*wordTokenList).InsertWordToken(WordToken.New("Andrea", 5))

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

	(*data).InsertWordToken(WordToken.New("Andrea", 5))
	(*data).InsertWordToken(WordToken.New("Graziani", 5))
	(*data).InsertWordToken(WordToken.New("Diana", 5))

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

	var currentWordToken *WordToken.WordToken

	wordTokenList1 := New()
	wordTokenList2 := New()

	(*wordTokenList1).InsertWordToken(WordToken.New("Andrea", 5))
	(*wordTokenList1).InsertWordToken(WordToken.New("Graziani", 30))
	(*wordTokenList1).InsertWordToken(WordToken.New("Yumi", 26))
	(*wordTokenList2).InsertWordToken(WordToken.New("Andrea", 5))
	(*wordTokenList2).InsertWordToken(WordToken.New("Graziani", 30))
	(*wordTokenList2).InsertWordToken(WordToken.New("Yumi", 26))

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
