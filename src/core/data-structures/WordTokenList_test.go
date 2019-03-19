package data_structures

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_WordTokenList_Inserting(t *testing.T) {

	var currentWordToken *WordToken

	wordTokenList := BuildWordTokenList()

	(*wordTokenList).InsertWordToken(BuildWordToken("Andrea", 5))
	(*wordTokenList).InsertWordToken(BuildWordToken("Graziani", 30))
	(*wordTokenList).InsertWordToken(BuildWordToken("Yumi", 26))
	(*wordTokenList).InsertWordToken(BuildWordToken("Yumi", 4))
	(*wordTokenList).InsertWordToken(BuildWordToken("Andrea", 5))

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

func Test_WordTokenList_Serialization(t *testing.T) {

	wordTokenList := BuildWordTokenList()

	(*wordTokenList).InsertWordToken(BuildWordToken("Andrea", 5))
	(*wordTokenList).InsertWordToken(BuildWordToken("Graziani", 5))
	(*wordTokenList).InsertWordToken(BuildWordToken("Diana", 5))

	output := (*wordTokenList).Serialize()
	fmt.Println(output)
}

func Test_WordTokenList_WordTokenListMerging(t *testing.T) {

	var currentWordToken *WordToken

	wordTokenList1 := BuildWordTokenList()
	wordTokenList2 := BuildWordTokenList()

	(*wordTokenList1).InsertWordToken(BuildWordToken("Andrea", 5))
	(*wordTokenList1).InsertWordToken(BuildWordToken("Graziani", 30))
	(*wordTokenList1).InsertWordToken(BuildWordToken("Yumi", 26))
	(*wordTokenList2).InsertWordToken(BuildWordToken("Andrea", 5))
	(*wordTokenList2).InsertWordToken(BuildWordToken("Graziani", 30))
	(*wordTokenList2).InsertWordToken(BuildWordToken("Yumi", 26))

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
