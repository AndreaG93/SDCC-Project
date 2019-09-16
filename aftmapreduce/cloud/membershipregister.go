package cloud

import (
	"fmt"
)

type MembershipRegister struct {
	register map[int]map[int]string
}

func (obj *MembershipRegister) GetWorkerProcessPublicInternetAddressesForRPC(groupId int, rpcBasePort int) []string {

	output := make([]string, len((*obj).register[groupId]))

	index := 0
	for nodeId, publicInternetAddress := range (*obj).register[groupId] {
		output[index] = fmt.Sprintf("%s:%d", publicInternetAddress, nodeId+rpcBasePort)
		index++
	}

	return output
}

func (obj *MembershipRegister) GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(groupId int, processIds []int, rpcBasePort int) []string {

	output := make([]string, len(processIds))

	for index, processId := range processIds {

		publicInternetAddress := (*obj).register[groupId][processId]
		output[index] = fmt.Sprintf("%s:%d", publicInternetAddress, processId+rpcBasePort)
	}

	return output
}

func (obj *MembershipRegister) GetWorkerProcessIDs(groupId int) []int {

	output := make([]int, len((*obj).register[groupId]))

	index := 0
	for nodeId, _ := range (*obj).register[groupId] {
		output[index] = nodeId
		index++
	}
	return output
}
