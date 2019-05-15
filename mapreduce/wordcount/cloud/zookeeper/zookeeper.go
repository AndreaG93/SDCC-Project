package zookeeper

import (
	"SDCC-Project-WorkerNode/utility"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

const (
	PrimaryNode         = "/primarynode"
	leaderZookeeperPath = "/leader"
)

func Init() {

	var zookeeperConnection *zk.Conn
	var isPathExisting bool
	var err error

	zookeeperConnection, _, err = zk.Connect([]string{"localhost"}, time.Second)
	utility.CheckError(err)

	isPathExisting, _, err = zookeeperConnection.Exists(leaderZookeeperPath)
	utility.CheckError(err)

	if !isPathExisting {

		_, err := zookeeperConnection.Create(leaderZookeeperPath, nil, 0, zk.WorldACL(zk.PermAll))
		utility.CheckError(err)

	}
}

func SetActualClusterLeaderAddress(address string) {

	var zookeeperConnection *zk.Conn
	var actualStat *zk.Stat
	var err error

	zookeeperConnection, _, err = zk.Connect([]string{"localhost"}, time.Second)
	utility.CheckError(err)

	_, actualStat, err = zookeeperConnection.Get(leaderZookeeperPath)
	utility.CheckError(err)

	_, err = zookeeperConnection.Set(leaderZookeeperPath, []byte(address), actualStat.Version)
	utility.CheckError(err)
}

func GetActualClusterLeaderAddress() (string, <-chan zk.Event) {

	var zookeeperConnection *zk.Conn
	var outputWatchEvent <-chan zk.Event
	var output []byte
	var err error

	zookeeperConnection, _, err = zk.Connect([]string{"localhost"}, time.Second)
	utility.CheckError(err)

	output, _, outputWatchEvent, err = zookeeperConnection.GetW(leaderZookeeperPath)
	utility.CheckError(err)

	return string(output), outputWatchEvent
}
