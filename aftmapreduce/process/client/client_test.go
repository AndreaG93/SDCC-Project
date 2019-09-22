package client

import "testing"

func TestClient(t *testing.T) {
	StartWork("../../../test-input-data/input1.txt", []string{"localhost:2181"})
}
