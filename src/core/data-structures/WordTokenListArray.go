package data_structures

type WordTokenListArraySerializedUnit struct {
}

type WordTokenListArray []*WordTokenList
type WordTokenListArraySerialized []WordToken

func (obj *WordTokenListArray) Built(size uint) {
	*obj = make(WordTokenListArray, size)
}

func BuildWordTokenListArray(size uint) *WordTokenListArray {
	return nil
}

func ReadWordTokenListArrayFromLocalFile(filename string) (WordTokenListArray, error) {

	return nil, nil
}
