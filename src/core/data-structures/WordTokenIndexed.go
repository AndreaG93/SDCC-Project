package data_structures

type WordTokenIndexedArray []WordTokenIndexed

type WordTokenIndexed struct {
	TableSizeOrIndex uint
	Word             string
	Occurrences      uint
}

func (obj WordTokenIndexedArray) Deserialize() *WordTokenHashTable {

	var output *WordTokenHashTable
	var currentWordTokenList *WordTokenList

	output = BuildWordTokenHashTable(obj[0].TableSizeOrIndex)

	for index := uint(1); index < uint(len(obj)); index++ {

		currentWord := input[index].Word
		currentWordOccurrences := input[index].Occurrences

		currentWordToken = BuildWordToken(currentWord, currentWordOccurrences)

		output.InsertWordToken(currentWordToken)
	}

	return output

	return nil
}
