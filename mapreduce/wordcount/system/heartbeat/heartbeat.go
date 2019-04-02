package heartbeat

import (
	"fmt"
	"github.com/codeskyblue/heartbeat"
	"net/http"
	"time"
)

func StartToSendHeartbeat(clientIdentifier string, serverAddress string) {
	client := &heartbeat.Client{
		ServerAddr: serverAddress,
		Secret:     "my-secret",
		Identifier: clientIdentifier,
	}
	cancel := client.Beat(500 * time.Millisecond)
	//defer cancel() // cancel heartbeat
	// Do something else

	if clientIdentifier == "Worker 3" {
		defer cancel()
	}
}

func StartToReceiveHeartbeat() {
	hbs := heartbeat.NewServer("my-secret", 3*time.Second) // secret: my-secret, timeout: 15s
	hbs.OnReconnect = func(identifier string, req *http.Request) {
		fmt.Println(identifier, "ssss")
	}

	hbs.OnConnect = func(identifier string, r *http.Request) {
		fmt.Println(identifier, "is online")
	}
	hbs.OnDisconnect = func(identifier string) {
		fmt.Println(identifier, "is offline")
	}
	http.Handle("/heartbeat", hbs)
	http.ListenAndServe(":7000", nil)
}
