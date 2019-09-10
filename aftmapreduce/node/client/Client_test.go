package client

import "testing"

func TestClient(t *testing.T) {
	StartWork("../../../test-input-data/input1.txt", []string{"3.89.215.191:2181", "3.89.215.191:2181", "18.212.39.96"})
}
