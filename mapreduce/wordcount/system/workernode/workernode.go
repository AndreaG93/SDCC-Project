package workernode

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount"
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
)

type WorkerNode struct {
	name             string
	listenPortForRPC string
}

func New(primaryNodeName string) *WorkerNode {

	output := new(WorkerNode)

	output.name = primaryNodeName
	return output
}

func (obj *WorkerNode) StartToRespondToRPCRequests() {

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Map{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := system.StartAcceptingRPCRequest(wordcount.Reduce{}, (*obj).listenPortForRPC); err != nil {
			panic(err)
		}
	}()
}
