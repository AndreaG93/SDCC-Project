package wordcount

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/node"
)

func startMapTask(splits []string, faultToleranceLevel int, workersInternetAddress map[int]string) {

	for groupIndex, split := range splits {

		internetAddresses := node.GetZookeeperClient().GetWorkerInternetAddressesForRPC(groupIndex, aftmapreduce.MapTaskRPCBasePort)

		go startArbitraryFaultTolerantMapTask(split, internetAddresses)
	}
}

func startArbitraryFaultTolerantMapTask(s string, strings map[int]string) {

}
