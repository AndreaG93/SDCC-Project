package zookeeper

import (
	"SDCC-Project-WorkerNode/utility"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"testing"
)

func yyTestSetCurrentMasterIPAddress(t *testing.T) {

	const Data = "Test"

	var output string
	var err error

	if err = SetCurrentMasterIPAddress(Data); err != nil {
		panic(err)
	}

	if output, err = GetCurrentMasterIPAddress(); err != nil {
		panic(err)
	}

	if strings.Compare(Data, output) != 0 {
		panic(err)
	}

}

func TestPrimaryAddressChange(t *testing.T) {

	var zkLeaderChangeEventChannel <-chan zk.Event
	var zkConnection *zk.Conn
	var data []byte
	var err error

	zkConnection, _, err = connectToZookeeperServers([]string{"localhost"})
	utility.CheckError(err)

	//checkExistenceOrGenerateZNode(zkConnection, "/current_leader_id", 0)

	for {

		data, _, zkLeaderChangeEventChannel, err = zkConnection.GetW("/current_leader_id")
		utility.CheckError(err)

		fmt.Println(string(data))

		<-zkLeaderChangeEventChannel
	}
}
