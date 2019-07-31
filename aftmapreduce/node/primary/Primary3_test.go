package primary

import "testing"

func Test_primary3(t *testing.T) {
	New(3, "127.0.0.1", []string{"127.0.0.1"}).StartWork()
}
