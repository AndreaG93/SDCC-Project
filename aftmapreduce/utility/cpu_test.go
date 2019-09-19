package utility

import (
	"fmt"
	"testing"
)

func Test_CPUPercentageUtilization(t *testing.T) {

	if output, err := GetCPUPercentageUtilizationAsInteger(); err != nil {
		panic(err)
	} else {
		fmt.Println(output)
	}
}
