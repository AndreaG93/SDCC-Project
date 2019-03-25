package wordtokenlist

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/datastructures/wordtoken"
	"SDCC-Project-WorkerNode/utility"
	"container/list"
	"fmt"
	"strings"
)

type WordTokenList struct {
	wordTokenList           *list.List
	currentWordTokenElement *list.Element
	length                  uint
}

func New() *WordTokenList {

	output := new(WordTokenList)

	(*output).wordTokenList = list.New()
	(*output).currentWordTokenElement = nil
	(*output).length = 0

	return output
}

func Deserialize(input []byte) (*WordTokenList, error) {

	output := New()
	serializedData := []wordtoken.WordToken{}

	if err := utility.Decode(input, &serializedData); err != nil {
		return nil, err
	}

	for index := uint(0); index < uint(len(serializedData)); index++ {
		(*output).InsertWordToken(&serializedData[index])
	}

	return output, nil
}

func (obj *WordTokenList) InsertWord(word string) {
	(*obj).InsertWordToken(wordtoken.New(word, 1))
}

func (obj *WordTokenList) InsertWordToken(wordToken *wordtoken.WordToken) {

	wordTokenList := (*obj).wordTokenList

	for e := (*wordTokenList).Front(); e != nil; e = (*e).Next() {

		currentWordToken := (*e).Value.(*wordtoken.WordToken)

		if strings.Compare((*currentWordToken).Word, (*wordToken).Word) == 0 {

			(*currentWordToken).Occurrences += (*wordToken).Occurrences

			return

		} else if strings.Compare((*currentWordToken).Word, (*wordToken).Word) > 0 {

			(*wordTokenList).InsertBefore(wordToken, e)
			(*obj).length++

			return
		}
	}

	(*wordTokenList).PushBack(wordToken)
	(*obj).length++
	return
}

func (obj *WordTokenList) Print() {

	wordTokenList := (*obj).wordTokenList

	for e := (*wordTokenList).Front(); e != nil; e = (*e).Next() {

		currentWordToken := e.Value.(*wordtoken.WordToken)

		fmt.Println(*currentWordToken)
	}
}

func (obj *WordTokenList) WordToken() *wordtoken.WordToken {

	var output *wordtoken.WordToken

	currentWordTokenElement := (*obj).currentWordTokenElement
	output = ((*currentWordTokenElement).Value).(*wordtoken.WordToken)

	return output
}

func (obj *WordTokenList) IteratorReset() {

	(*obj).currentWordTokenElement = nil
}

func (obj *WordTokenList) Next() bool {

	if (*obj).currentWordTokenElement != nil {

		(*obj).currentWordTokenElement = (*(*obj).currentWordTokenElement).Next()

		if (*obj).currentWordTokenElement != nil {
			return true
		} else {
			return false
		}
	} else {

		wordTokenList := (*obj).wordTokenList
		(*obj).currentWordTokenElement = wordTokenList.Front()

		if (*obj).currentWordTokenElement != nil {
			return true
		} else {
			return false
		}
	}
}

func (obj *WordTokenList) Merge(input *WordTokenList) {

	(*obj).IteratorReset()

	for input.Next() {

		currentWordToken := input.WordToken()
		(*obj).InsertWordToken(currentWordToken)
	}
}

func (obj *WordTokenList) Serialize() ([]byte, error) {

	(*obj).IteratorReset()

	output := make([]wordtoken.WordToken, (*obj).length)

	for index := uint(0); index < (*obj).length; index++ {

		(*obj).Next()

		currentWordToken := (*obj).WordToken()

		output[index].Word = (*currentWordToken).Word
		output[index].Occurrences = (*currentWordToken).Occurrences
	}

	return utility.Encode(output)
}

func (obj *WordTokenList) GetLength() uint {
	return (*obj).length
}
