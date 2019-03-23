package wordtoken

type WordToken struct {
	Word        string
	Occurrences uint
}

func New(word string, numberOfOccurrences uint) *WordToken {

	output := new(WordToken)

	(*output).Word = word
	(*output).Occurrences = numberOfOccurrences

	return output
}
