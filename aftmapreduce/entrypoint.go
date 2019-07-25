package aftmapreduce

import "SDCC-Project/aftmapreduce/data"

type EntryPoint struct {
}

type EntryPointInput struct {
	Data data.ClientData
}

type EntryPointOutput struct {
}

func (x *EntryPoint) Execute(input EntryPointInput, output *EntryPointOutput) error {

	clientRequest := NewRequest(&input.Data)
	go ManageClientRequest(clientRequest)

	return nil
}
