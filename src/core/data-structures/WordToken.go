package data_structures

type WordToken struct {
	Word        string
	Occurrences uint
}

func BuildWordToken(word string) *WordToken {

	output := new(WordToken)
	output.Word = word
	output.Occurrences = 1

	return output
}
