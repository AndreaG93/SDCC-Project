package Primary

import "testing"

func Test_primary(t *testing.T) {
	New(1, "127.0.0.1:15001").StartWork()
}
