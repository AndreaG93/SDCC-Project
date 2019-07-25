package aftmapreduce

type EntryPoint struct {
}

type EntryPointInput struct {
	Data ClientData
}

type EntryPointOutput struct {
}

func (x *EntryPoint) Execute(input EntryPointInput, output *EntryPointOutput) error {

	clientRequest := NewRequest(&input.Data)
	go ManageClientRequest(clientRequest)

	return nil
}
