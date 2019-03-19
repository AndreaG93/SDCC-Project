package data_structures

type WordToken struct {
	Word        string
	Occurrences uint
}

func BuildWordToken(word string, numberOfOccurrences uint) *WordToken {

	output := new(WordToken)
	output.Word = word
	output.Occurrences = numberOfOccurrences

	return output
}

func BuildWordTokenByWord(word string) *WordToken {
	return BuildWordToken(word, 1)
}
