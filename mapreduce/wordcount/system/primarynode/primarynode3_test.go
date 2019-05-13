package primarynode

import (
	"testing"
)

func TestPrimaryNode3(t *testing.T) {

	primaryNode1 := New(3)
	(*primaryNode1).StartWork()
}
