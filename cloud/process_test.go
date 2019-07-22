package cloud

import (
	"SDCC-Project/utility"
	"fmt"
	"github.com/codeskyblue/heartbeat"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func Test_2(t *testing.T) {

	cmd := exec.Command("go", "run", "../main/HeartbeatSenderProcess.go", "localhost", "10000", "/heartbeat", "1")
	utility.CheckError(cmd.Start())

	startReceivingHeartbeat()

}

func startReceivingHeartbeat() {

	hbs := heartbeat.NewServer("my-secret", 15*time.Second) // secret: my-secret, timeout: 15s
	hbs.OnConnect = func(identifier string, r *http.Request) {
		fmt.Println(identifier, "is online")
	}
	hbs.OnDisconnect = func(identifier string) {
		fmt.Println(identifier, "is offline")
	}
	http.Handle("/heartbeat", hbs)
	http.ListenAndServe(":10000", nil)
}
