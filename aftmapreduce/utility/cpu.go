package utility

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func GetCPUPercentageUtilizationAsInteger() (int, error) {

	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		command = exec.Command("wsl.exe", "/bin/bash", "-c", "mpstat | grep -A 5 \"%idle\" | tail -n 1 | awk -F \" \" '{print 100 -  $ 12}'a")
	} else {
		command = exec.Command("/bin/bash", "-c", "mpstat | grep -A 5 \"%idle\" | tail -n 1 | awk -F \" \" '{print 100 -  $ 12}'a")
	}

	if output, err := command.Output(); err != nil {
		return 0, err
	} else {
		if outputAsFloat, err := strconv.ParseFloat(strings.TrimSuffix(string(output), "\n"), 64); err != nil {
			return 0, err
		} else {
			return int(outputAsFloat), nil
		}
	}
}
