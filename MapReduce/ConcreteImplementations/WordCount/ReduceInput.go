package WordCount

import (
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenListGroup"
	"SDCC-Project/utility"
)

type ReduceInput struct {
	Data []byte
}

func (obj *ReduceInput) PerformTask() (string, []byte, error) {

	var err error

	tokenListGroup, err := WordTokenListGroup.Deserialize((*obj).Data)
	if err != nil {
		return "", nil, err
	}

	data := tokenListGroup.Merge()

	rawData, err := data.Serialize()

	digest := utility.GenerateDigestUsingSHA512(rawData)

	return digest, rawData, nil
}
