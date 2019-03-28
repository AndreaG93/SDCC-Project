package zookeeper

import (
	"strings"
	"testing"
)

func TestSetCurrentMasterIPAddress(t *testing.T) {

	const Data = "Test"

	var output string
	var err error

	if err = SetCurrentMasterIPAddress(Data); err != nil {
		panic(err)
	}

	if output, err = GetCurrentMasterIPAddress(); err != nil {
		panic(err)
	}

	if strings.Compare(Data, output) != 0 {
		panic(err)
	}

}
