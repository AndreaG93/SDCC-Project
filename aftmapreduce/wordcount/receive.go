package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
)

type Receive struct {
}

type ReceiveInput struct {
	Data               *WordTokenList.WordTokenList
	ReceivedDataDigest string
}

type ReceiveOutput struct {
}

func (x *Receive) Execute(input ReceiveInput, output *ReceiveOutput) error {

	node.GetCache().Set(input.ReceivedDataDigest, input.Data)
	return nil
}
