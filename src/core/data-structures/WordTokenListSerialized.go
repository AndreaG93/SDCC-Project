package data_structures

type WordTokenListSerialized []WordToken

func (obj WordTokenListSerialized) Deserialize() *WordTokenList {

	output := BuildWordTokenList()

	for index := uint(0); index < uint(len(obj)); index++ {
		(*output).InsertWordToken(&obj[index])
	}

	return output
}
