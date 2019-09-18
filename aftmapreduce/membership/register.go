package membership

import (
	"fmt"
)

type Register struct {
	coordinator Coordinator
	register    map[int]map[int]string
}

func New(coordinator Coordinator) (*Register, error) {

	var err error
	output := new(Register)
	(*output).coordinator = coordinator

	if (*output).register, err = ((*output).coordinator).GetProcessMembershipTable(); err != nil {
		return nil, err
	} else {
		go (*output).startListeningForMembershipChanges()
		return output, nil
	}
}

func (obj *Register) startListeningForMembershipChanges() {

	var err error

	for {

		if (*obj).register, err = ((*obj).coordinator).GetProcessMembershipTable(); err != nil {
			panic(err)
		}

		if err = ((*obj).coordinator).WaitUntilProcessMembershipChanges(); err != nil {
			panic(err)
		}
	}
}

func (obj *Register) GetWorkerProcessPublicInternetAddressesForRPC(groupId int, rpcBasePort int) ([]string, error) {

	output := make([]string, len((*obj).register[groupId]))

	index := 0
	for nodeId, publicInternetAddress := range (*obj).register[groupId] {
		output[index] = fmt.Sprintf("%s:%d", publicInternetAddress, nodeId+rpcBasePort)
		index++
	}

	return networkLocalitySort(output)
}

func (obj *Register) GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(groupId int, processIds []int, rpcBasePort int) ([]string, error) {

	output := make([]string, len(processIds))

	for index, processId := range processIds {

		publicInternetAddress := (*obj).register[groupId][processId]
		output[index] = fmt.Sprintf("%s:%d", publicInternetAddress, processId+rpcBasePort)
	}

	return networkLocalitySort(output)
}

func (obj *Register) GetWorkerProcessIDs(groupId int) []int {

	output := make([]int, len((*obj).register[groupId]))

	index := 0
	for nodeId := range (*obj).register[groupId] {
		output[index] = nodeId
		index++
	}
	return output
}

func (obj *Register) GetGroupAmount() int {
	return len((*obj).register)
}
