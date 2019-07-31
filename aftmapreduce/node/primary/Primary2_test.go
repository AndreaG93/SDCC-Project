package primary

import "testing"

func Test_primary2(t *testing.T) {
	New(2, "127.0.0.1", []string{"127.0.0.1"}).StartWork()
}
