package primary

import "testing"

func Test_primary4(t *testing.T) {
	New(4, "127.0.0.1", []string{"127.0.0.1"}).StartWork()
}
