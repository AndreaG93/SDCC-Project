package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/registry"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"SDCC-Project/utility"
	"fmt"
)

type Receive struct {
}

type ReceiveInput struct {
	Data                 []byte
	ReceivedDataDigest   string
	AssociatedDataDigest string
}

type ReceiveOutput struct {
}

func (x *Receive) Execute(input ReceiveInput, output *ReceiveOutput) error {

	node.GetLogger().PrintMessage(fmt.Sprintf("Received Data Digest: %s Associated to Data Digest: %s", input.ReceivedDataDigest, input.AssociatedDataDigest))

	receivedWordTolenList, err := WordTokenList.Deserialize(input.Data)
	utility.CheckError(err)

	node.GetCache().Set(input.ReceivedDataDigest, receivedWordTolenList)
	registry.GetDigestCacheInstance().Add(input.AssociatedDataDigest, input.ReceivedDataDigest)

	return nil
}
