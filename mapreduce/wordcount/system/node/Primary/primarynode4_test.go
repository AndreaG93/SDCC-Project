package Primary

import (
	"testing"
)

func TestPrimaryNode4(t *testing.T) {

	primaryNode1 := New(4)
	(*primaryNode1).StartWork()
}
