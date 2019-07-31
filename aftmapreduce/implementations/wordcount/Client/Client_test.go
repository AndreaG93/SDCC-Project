package Client

import "testing"

func TestClient(t *testing.T) {
	StartWork("../../../../test-input-data/input1.txt", []string{"127.0.0.1:2181"})
}
