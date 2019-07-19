package zookeeper

import (
	"SDCC-Project-WorkerNode/utility"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

const (
	ActualLeaderZNodePath = "/leader"
	lockerTimeout         = 1 * time.Minute
	zkSessionTimeOut      = 20 * time.Second
)

type Client struct {
	zooKeeperConnection *zk.Conn
	zooKeeperLock       *zk.Lock
}

func New(zooKeeperServerPoolAddresses []string) *Client {

	var err error
	output := new(Client)

	(*output).zooKeeperConnection, _, err = zk.Connect(zooKeeperServerPoolAddresses, zkSessionTimeOut)
	(*output).zooKeeperLock = nil
	utility.CheckError(err)

	return output
}

func (obj *Client) CheckZNodeExistence(zNodePath string) bool {

	var output bool
	var err error

	output, _, err = (*obj).zooKeeperConnection.Exists(zNodePath)
	utility.CheckError(err)

	return output
}

func (obj *Client) CreateZNode(zNodePath string) {

	_, err := (*obj).zooKeeperConnection.Create(zNodePath, nil, 0, zk.WorldACL(zk.PermAll))
	utility.CheckError(err)

}

func (obj *Client) RemoveZNode(zNodePath string) {

	var zNodeExistence bool
	var actualStat *zk.Stat
	var err error

	zNodeExistence, actualStat, err = (*obj).zooKeeperConnection.Exists(zNodePath)
	utility.CheckError(err)

	if zNodeExistence {

		err := (*obj).zooKeeperConnection.Delete(zNodePath, actualStat.Version)
		utility.CheckError(err)

	}
}

func (obj *Client) SetZNodeData(zNodePath string, data []byte) {

	var actualStat *zk.Stat
	var err error

	_, actualStat, err = (*obj).zooKeeperConnection.Get(zNodePath)
	utility.CheckError(err)

	_, err = (*obj).zooKeeperConnection.Set(zNodePath, data, actualStat.Version)
	utility.CheckError(err)

}

func (obj *Client) GetZNodeData(zNodePath string) ([]byte, <-chan zk.Event) {

	outputData, _, outputWatchEvent, err := (*obj).zooKeeperConnection.GetW(zNodePath)
	utility.CheckError(err)

	return outputData, outputWatchEvent
}

func (obj *Client) CloseConnection() {
	(*obj).zooKeeperConnection.Close()
}
