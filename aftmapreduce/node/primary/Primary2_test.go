package primary

import "testing"

func Test_primary2(t *testing.T) {
	Initialize(2, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
