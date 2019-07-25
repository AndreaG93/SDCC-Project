package Client

import "testing"

func TestClient(t *testing.T) {
	SendRequest("../../../test-input-data/input1.txt", "127.0.0.1")
}
