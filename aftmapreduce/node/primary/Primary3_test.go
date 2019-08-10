package primary

import "testing"

func Test_primary3(t *testing.T) {
	Initialize(3, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
