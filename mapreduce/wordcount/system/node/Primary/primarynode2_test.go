package Primary

import (
	"testing"
)

func TestPrimaryNode2(t *testing.T) {

	primaryNode1 := New(2)
	(*primaryNode1).StartWork()
}
