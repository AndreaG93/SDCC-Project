package membership

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func networkLocalitySort(publicInternetAddresses []string) ([]string, error) {

	output := make([]string, len(publicInternetAddresses))
	averageResponseTimes := make([]float64, 0)
	mapStructure := make(map[float64]string)

	for _, ip := range publicInternetAddresses {

		if averageResponseTime, err := getRoundTripTime(strings.Split(ip, ":")[0]); err != nil {
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

func getRoundTripTime(publicInternetAddress string) (float64, error) {

	var command *exec.Cmd
	var err error
	var averageResponseTimeString string
	var averageResponseTime float64

	if runtime.GOOS == "windows" {
		command = exec.Command("wsl.exe", "/bin/bash", "-c", fmt.Sprintf("ping -c 1 %s | cut -d '/' -s -f5", publicInternetAddress))
	} else {
		command = exec.Command("/bin/bash", "-c", fmt.Sprintf("ping -c 1 %s | cut -d '/' -s -f5", publicInternetAddress))
	}

	output, err := command.Output()
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
