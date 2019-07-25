package aftmapreduce

import (
	"SDCC-Project/aftmapreduce/data"
	"SDCC-Project/aftmapreduce/registries/WorkerMutex"
	"SDCC-Project/aftmapreduce/registries/WorkerResultsRegister"
)

type Replica struct {
}

type ReplicaInput struct {
	Data data.TransientData
}

type ReplicaOutput struct {
	Digest string
}

func (x *Replica) Execute(input ReplicaInput, output *ReplicaOutput) error {

	digest, rawData, err := input.Data.PerformTask()
	if err != nil {
		return err
	}

	WorkerMutex.GetInstance().Lock()
	WorkerResultsRegister.GetInstance().Set(digest, rawData)
	WorkerMutex.GetInstance().Unlock()

	output.Digest = digest

	return nil
}
