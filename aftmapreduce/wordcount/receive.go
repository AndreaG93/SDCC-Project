package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"fmt"
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
	fmt.Println(input.Data)
	node.GetCache().Set(input.ReceivedDataDigest, input.Data)
	return nil
}
