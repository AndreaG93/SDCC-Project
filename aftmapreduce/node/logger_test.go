package node

import "testing"

func Test_worker1(t *testing.T) {
	logger := New(1, "PrimaryNode")

	logger.PrintMessage("Test")
}
