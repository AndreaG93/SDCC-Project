package zookeeper

import (
	"SDCC-Project/utility"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

const (
	membershipZNodeRootPath = "/membership"
	ActualLeaderZNodePath   = "/leader"
	zkSessionTimeOut        = 20 * time.Second
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

	if !(*output).CheckZNodeExistence(membershipZNodeRootPath) {
		(*output).CreateZNode(membershipZNodeRootPath, 0)
	}

	return output
}

func (obj *Client) CheckZNodeExistence(zNodePath string) bool {

	var output bool
	var err error

	output, _, err = (*obj).zooKeeperConnection.Exists(zNodePath)
	utility.CheckError(err)

	return output
}

func (obj *Client) CreateZNode(zNodePath string, flags int32) {

	_, err := (*obj).zooKeeperConnection.Create(zNodePath, nil, flags, zk.WorldACL(zk.PermAll))
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

func (obj *Client) GetMembersInternetAddress() map[int]string {

	members, _ := (*obj).GetMembersList()
	membersInternetAddress := make(map[int]string)

	for _, element := range members {

		path := fmt.Sprintf("%s/%s", membershipZNodeRootPath, element)
		rawData, _ := (*obj).GetZNodeData(path)
		memberID, err := strconv.Atoi(element)
		utility.CheckError(err)

		membersInternetAddress[memberID] = string(rawData)
	}

	return membersInternetAddress
}

func (obj *Client) GetMembersList() ([]string, <-chan zk.Event) {

	data, _, channel, err := (*obj).zooKeeperConnection.ChildrenW(membershipZNodeRootPath)
	utility.CheckError(err)
	return data, channel
}

func (obj *Client) RegisterNodeMembership(nodeID int, internetAddress string) {

	path := fmt.Sprintf("%s/%d", membershipZNodeRootPath, nodeID)

	if !(*obj).CheckZNodeExistence(path) {
		(*obj).CreateZNode(path, zk.FlagEphemeral)
	}

	(*obj).SetZNodeData(path, []byte(internetAddress))
}

func (obj *Client) KeepConnectionAlive() {

	_, channel := (*obj).GetZNodeData("/")
	<-channel
}
