package Client

import "testing"

func TestClient(t *testing.T) {
	SendRequest("../../../test-input-data/input.txt", "127.0.0.1:15001")
}
