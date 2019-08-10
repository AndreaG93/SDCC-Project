package primary

import "testing"

func Test_primary1(t *testing.T) {
	Initialize(1, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
