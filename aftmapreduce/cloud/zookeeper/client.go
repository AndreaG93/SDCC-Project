package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

const (
	electionZNodePath                = "/election"
	membershipZNodePath              = "/membership"
	pendingClientRequestsZNodePath   = "/pending-requests"
	finalizedClientRequestsZNodePath = "/finalized-requests"
	zkSessionTimeOut                 = 15 * time.Second
)

type Client struct {
	zooKeeperConnection *zk.Conn
}

func New(zooKeeperServerPoolAddresses []string) (*Client, error) {

	var err error

	output := new(Client)
	(*output).zooKeeperConnection, _, err = zk.Connect(zooKeeperServerPoolAddresses, zkSessionTimeOut)

	go (*output).KeepConnectionAlive()

	return output, err
}

func (obj *Client) initializeAllNeededZNodes() error {

	var err error

	if err = (*obj).createZNodeCheckingExistence(electionZNodePath, nil, int32(0)); err != nil {
		return err
	}
	if err = (*obj).createZNodeCheckingExistence(membershipZNodePath, nil, int32(0)); err != nil {
		return err
	}
	if err = (*obj).createZNodeCheckingExistence(pendingClientRequestsZNodePath, nil, int32(0)); err != nil {
		return err
	}
	if err = (*obj).createZNodeCheckingExistence(finalizedClientRequestsZNodePath, nil, int32(0)); err != nil {
		return err
	}

	return nil
}

func (obj *Client) createZNodeCheckingExistence(zNodePath string, data []byte, flags int32) error {

	var err error
	var isZNodeExistent bool

	if isZNodeExistent, _, err = (*obj).zooKeeperConnection.Exists(electionZNodePath); err != nil {
		return err
	} else if !isZNodeExistent {
		if _, err = (*obj).zooKeeperConnection.Create(zNodePath, data, flags, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	return nil
}

func (obj *Client) removeZNodeCheckingExistence(zNodePath string) error {

	var err error
	var actualStat *zk.Stat
	var isZNodeExistent bool

	if isZNodeExistent, actualStat, err = (*obj).zooKeeperConnection.Exists(zNodePath); err != nil {
		return err
	} else if isZNodeExistent {
		return (*obj).zooKeeperConnection.Delete(zNodePath, actualStat.Version)
	} else {
		return nil
	}
}

func (obj *Client) getChildrenZNodeNames(zNodePath string) ([]string, <-chan zk.Event, error) {

	data, _, channel, err := (*obj).zooKeeperConnection.ChildrenW(zNodePath)
	return data, channel, err
}

func (obj *Client) setZNodeData(zNodePath string, data []byte) error {

	var err error

	if _, actualStat, err := (*obj).zooKeeperConnection.Get(zNodePath); err == nil {
		_, err = (*obj).zooKeeperConnection.Set(zNodePath, data, actualStat.Version)
	}

	return err
}

func (obj *Client) getZNodeData(zNodePath string) ([]byte, <-chan zk.Event, error) {

	outputData, _, outputWatchEvent, err := (*obj).zooKeeperConnection.GetW(zNodePath)
	return outputData, outputWatchEvent, err
}

func (obj *Client) KeepConnectionAlive() {

	if _, channel, err := (*obj).getZNodeData("/"); err == nil {
		<-channel
	} else {
		panic(err)
	}
}
