package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

const (
	PrimaryNode = "/primary"
)

func SetCurrentMasterIPAddress(address string) error {

	var zookeeperConnection *zk.Conn
	var actualStat *zk.Stat
	var isPathExisting bool
	var err error

	if zookeeperConnection, _, err = zk.Connect([]string{"localhost"}, time.Second); err != nil {
		return err
	}

	if isPathExisting, _, err = zookeeperConnection.Exists(PrimaryNode); err != nil {
		return err
	}

	if !isPathExisting {

		if _, err := zookeeperConnection.Create(PrimaryNode, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}

		if _, err := zookeeperConnection.Set(PrimaryNode, []byte(address), 0); err != nil {
			return err
		}

	} else {

		if _, actualStat, err = zookeeperConnection.Get(PrimaryNode); err != nil {
			return err
		}

		if _, err := zookeeperConnection.Set(PrimaryNode, []byte(address), actualStat.Version); err != nil {
			return err
		}

	}

	return nil
}

func GetCurrentMasterIPAddress() (string, error) {

	var zookeeperConnection *zk.Conn
	var isPathExisting bool
	var rawOutput []byte
	var err error

	if zookeeperConnection, _, err = zk.Connect([]string{"localhost"}, time.Second); err != nil {
		return "", err
	}

	if isPathExisting, _, err = zookeeperConnection.Exists(PrimaryNode); err != nil || !isPathExisting {
		return "", err
	}

	if rawOutput, _, err = zookeeperConnection.Get(PrimaryNode); err != nil {
		return "", err
	}

	return string(rawOutput), nil
}

func GetLocalClusterWorkerPopulation() uint {
	return 10
}
