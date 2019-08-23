package zookeeper

import (
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

const (
	membershipZNodeRootPath = "/membership"
	zkSessionTimeOut        = 15 * time.Second
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
		(*output).CreateZNode(membershipZNodeRootPath, nil, 0)
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

func (obj *Client) CreateZNode(zNodePath string, data []byte, flags int32) {

	_, err := (*obj).zooKeeperConnection.Create(zNodePath, data, flags, zk.WorldACL(zk.PermAll))
	utility.CheckError(err)
}

func (obj *Client) CreateZNodeCheckingExistence(zNodePath string, data []byte, flags int32) {

	output, _, err := (*obj).zooKeeperConnection.Exists(zNodePath)
	utility.CheckError(err)

	if !output {
		_, err := (*obj).zooKeeperConnection.Create(zNodePath, data, flags, zk.WorldACL(zk.PermAll))
		utility.CheckError(err)
	}
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

func (obj *Client) GetZNodeWatcher(zNodePath string) <-chan zk.Event {

	_, output := (*obj).GetZNodeData(zNodePath)
	return output
}

func (obj *Client) CloseConnection() {
	(*obj).zooKeeperConnection.Close()
}

func (obj *Client) getChildrenZNode(parentZNode string) ([]string, <-chan zk.Event) {

	data, _, channel, err := (*obj).zooKeeperConnection.ChildrenW(parentZNode)
	utility.CheckError(err)
	return data, channel
}

func (obj *Client) KeepConnectionAlive() {

	_, channel := (*obj).GetZNodeData("/")
	<-channel
}

func (obj *Client) GetChildrenList(zNodePath string) []string {

	output, _, err := (*obj).zooKeeperConnection.Children(zNodePath)
	utility.CheckError(err)
	return output
}

func (obj *Client) RegisterNodeMembership(nodeId int, groupId int, internetAddress string) {

	groupPath := fmt.Sprintf("%s/%d", membershipZNodeRootPath, groupId)
	nodePath := fmt.Sprintf("%s/%d", groupPath, nodeId)

	if !(*obj).CheckZNodeExistence(groupPath) {
		(*obj).CreateZNode(groupPath, nil, 0)
		(*obj).CreateZNode(nodePath, []byte(internetAddress), zk.FlagEphemeral)
	} else {
		if (*obj).CheckZNodeExistence(nodePath) {
			(*obj).RemoveZNode(nodePath)
			(*obj).CreateZNode(nodePath, []byte(internetAddress), zk.FlagEphemeral)
		} else {
			(*obj).CreateZNode(nodePath, []byte(internetAddress), zk.FlagEphemeral)
		}
	}
}

func (obj *Client) GetMappedWorkerInternetAddressesForRPC(groupId int, baseRPCPort int) map[int]string {

	output := make(map[int]string)
	parentZNodePath := fmt.Sprintf("%s/%d", membershipZNodeRootPath, groupId)
	group, _ := (*obj).getChildrenZNode(parentZNodePath)

	for _, element := range group {

		nodeId, err := strconv.Atoi(element)
		utility.CheckError(err)

		rawInternetAddresses, _ := (*obj).GetZNodeData(fmt.Sprintf("%s/%s", parentZNodePath, element))
		internetAddresses := fmt.Sprintf("%s:%d", string(rawInternetAddresses), nodeId+baseRPCPort)

		output[nodeId] = internetAddresses
	}

	return output
}

func (obj *Client) GetWorkerInternetAddressesForRPC(groupId int, baseRPCPort int) []string {

	parentZNodePath := fmt.Sprintf("%s/%d", membershipZNodeRootPath, groupId)
	availableWorkerNodes, _ := (*obj).getChildrenZNode(parentZNodePath)

	output := make([]string, len(availableWorkerNodes))

	for index, element := range availableWorkerNodes {

		nodeId, err := strconv.Atoi(element)
		utility.CheckError(err)

		rawInternetAddresses, _ := (*obj).GetZNodeData(fmt.Sprintf("%s/%s", parentZNodePath, element))
		internetAddresses := fmt.Sprintf("%s:%d", string(rawInternetAddresses), nodeId+baseRPCPort)

		output[index] = internetAddresses
	}

	return output
}

func (obj *Client) GetWorkerInternetAddressesForRPCWithIdConstraints(groupId int, baseRPCPort int, requiredId []int) []string {

	output := make([]string, len(requiredId))
	mappedWorkerInternetAddresses := (*obj).GetMappedWorkerInternetAddressesForRPC(groupId, baseRPCPort)

	index := 0
	for _, id := range requiredId {
		output[index] = mappedWorkerInternetAddresses[id]
		index++
	}

	return output
}

func (obj *Client) GetGroupAmount() int {

	groups, _ := (*obj).getChildrenZNode(membershipZNodeRootPath)
	return len(groups)
}
