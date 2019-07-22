package Client

import "testing"

func TestClient(t *testing.T) {
	SendRequest("../../../test-input-data/input2.txt", "127.0.0.1:15001")
}
