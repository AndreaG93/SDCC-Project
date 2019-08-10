package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/data"
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/node/property"
	"SDCC-Project/aftmapreduce/wordcount/aft"
	"SDCC-Project/cloud/amazons3"
	"SDCC-Project/utility"
	"fmt"
)

type Request struct {
}

type RequestInput struct {
	SourceDataDigest string
}

type RequestOutput struct {
}

func (x *Request) Execute(input RequestInput, output *RequestOutput) error {

	go manageRequest(input.SourceDataDigest)
	return nil
}

func manageRequest(digest string) {

	splits := getSplits(digest, node.GetPropertyAsInteger(property.MapCardinality))

	for index, split := range splits {
		aft.NewMapTask()
	}

}
