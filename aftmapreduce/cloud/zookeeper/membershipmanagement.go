package zookeeper

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
)

func (obj *Client) RegisterNewWorkerProcess(processId int, processGroupId int, processPublicInternetAddress string) error {

	zNodePath := fmt.Sprintf("%s/%d_%d", membershipZNodePath, processGroupId, processId)
	return (*obj).createZNodeCheckingExistence(zNodePath, []byte(processPublicInternetAddress), zk.FlagEphemeral)
}

func (obj *Client) WaitUntilProcessMembershipChanges() error {

	var err error
	var watcher <-chan zk.Event

	if _, _, watcher, err = (*obj).zooKeeperConnection.ChildrenW(membershipZNodePath); err != nil {
		return err
	}

	<-watcher
	return nil
}

func (obj *Client) UpdateProcessMembershipRegister() (map[int]map[int]string, error) {

	var err error
	var zNodeNames []string

	if zNodeNames, _, err = (*obj).zooKeeperConnection.Children(membershipZNodePath); err != nil {
		return nil, err
	}

	return (*obj).extractProcessMembershipTableFrom(zNodeNames)
}

func (obj *Client) extractProcessMembershipTableFrom(zNodeNames []string) (map[int]map[int]string, error) {

	output := make(map[int]map[int]string)

	for _, zNodeName := range zNodeNames {

		groupID, nodeID, err := extractProcessInfoFrom(zNodeName)
		if err != nil {
			return nil, err
		}

		publicInternetAddress, _, err := (*obj).getZNodeData(fmt.Sprintf("%s/%s", membershipZNodePath, zNodeName))
		if err != nil {
			return nil, err
		}

		if output[groupID] == nil {
			output[groupID] = make(map[int]string)
		}

		output[groupID][nodeID] = string(publicInternetAddress)
	}

	return output, nil
}

func extractProcessInfoFrom(zNodeName string) (int, int, error) {

	var groupID int
	var nodeID int
	var err error

	raw := strings.Split(zNodeName, "_")

	if groupID, err = strconv.Atoi(raw[0]); err != nil {
		return -1, -1, err
	}
	if nodeID, err = strconv.Atoi(raw[1]); err != nil {
		return -1, -1, err
	}

	return groupID, nodeID, nil
}
