package Primary

import (
	"testing"
)

func TestPrimaryNode1(t *testing.T) {

	primaryNode1 := New(1)
	(*primaryNode1).StartWork()
}
