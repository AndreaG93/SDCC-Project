package cloud

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type MembershipRegister struct {
	register map[int]map[int]string
}

func NewMembershipRegister() *MembershipRegister {

	output := new(MembershipRegister)
	return output
}

func (obj *MembershipRegister) StartMembershipRegisterListener(systemCoordinator SystemCoordinator) {

	var err error

	for {

		if (*obj).register, err = systemCoordinator.UpdateProcessMembershipRegister(); err != nil {
			panic(err)
		}

		if err = systemCoordinator.WaitUntilProcessMembershipChanges(); err != nil {
			panic(err)
		}

	}
}

func (obj *MembershipRegister) GetWorkerProcessPublicInternetAddressesForRPC(groupId int, rpcBasePort int) ([]string, error) {

	output := make([]string, len((*obj).register[groupId]))

	index := 0
	for nodeId, publicInternetAddress := range (*obj).register[groupId] {
		output[index] = fmt.Sprintf("%s:%d", publicInternetAddress, nodeId+rpcBasePort)
		index++
	}

	return orderPublicInternetRPCAddressByAverageResponseTime(output)
}

func (obj *MembershipRegister) GetSpecifiedWorkerProcessPublicInternetAddressesForRPC(groupId int, processIds []int, rpcBasePort int) ([]string, error) {

	output := make([]string, len(processIds))

	for index, processId := range processIds {

		publicInternetAddress := (*obj).register[groupId][processId]
		output[index] = fmt.Sprintf("%s:%d", publicInternetAddress, processId+rpcBasePort)
	}

	return orderPublicInternetRPCAddressByAverageResponseTime(output)
}

func (obj *MembershipRegister) GetWorkerProcessIDs(groupId int) []int {

	output := make([]int, len((*obj).register[groupId]))

	index := 0
	for nodeId := range (*obj).register[groupId] {
		output[index] = nodeId
		index++
	}
	return output
}

func (obj *MembershipRegister) GetGroupAmount() int {
	return len((*obj).register)
}

func orderPublicInternetRPCAddressByAverageResponseTime(input []string) ([]string, error) {

	output := make([]string, len(input))
	averageResponseTimes := make([]float64, 0)
	mapStructure := make(map[float64]string)

	for _, ip := range input {

		if averageResponseTime, err := getAverageResponseTime(strings.Split(ip, ":")[0]); err != nil {
			return nil, err
		} else {

			if mapStructure[averageResponseTime] != "" {
				averageResponseTime += averageResponseTime / 1000000000
			}

			mapStructure[averageResponseTime] = ip
			averageResponseTimes = append(averageResponseTimes, averageResponseTime)
		}
	}

	sort.Float64s(averageResponseTimes)

	for index, x := range averageResponseTimes {
		output[index] = mapStructure[x]
	}

	return output, nil
}

func getAverageResponseTime(publicInternetAddress string) (float64, error) {

	var err error
	var averageResponseTimeString string
	var averageResponseTime float64

	cmd := exec.Command("wsl.exe", "/bin/bash", "-c", fmt.Sprintf("ping -c 1 %s | cut -d '/' -s -f5", publicInternetAddress))
	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}
	averageResponseTimeString = strings.TrimSuffix(string(output), "\n")
	averageResponseTime, err = strconv.ParseFloat(averageResponseTimeString, 64)
	if err != nil {
		return 0.0, err
	}

	return averageResponseTime, nil
}
