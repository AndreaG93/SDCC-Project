package node

import (
	"SDCC-Project/aftmapreduce/data"
	"SDCC-Project/aftmapreduce/membership"
	"SDCC-Project/aftmapreduce/storage"
	"SDCC-Project/aftmapreduce/storage/amazons3"
	"SDCC-Project/aftmapreduce/system"
	"SDCC-Project/aftmapreduce/system/zookeeper"
)

var systemCoordinator system.Coordinator
var keyValueRegister storage.KeyValueRegister
var membershipRegister *membership.Register
var membershipCoordinator membership.Coordinator
var dataRegistry *data.Registry

var properties map[string]interface{}

var logger *Logger

func InitializePrimary(guid uint, zookeeperAddresses []string) error {

	var err error

	if zookeeperClient, err := zookeeper.New(zookeeperAddresses); err != nil {
		return err
	} else {

		systemCoordinator = zookeeperClient
		membershipCoordinator = zookeeperClient
	}

	if err = systemCoordinator.Initialize(); err != nil {
		return err
	}
	if dataRegistry, err = data.New(guid, "Primary", true); err != nil {
		return err
	}
	if membershipRegister, err = membership.New(membershipCoordinator); err != nil {
		return err
	}

	keyValueRegister = amazons3.New()
	logger = NewLogger()
	properties = make(map[string]interface{})

	return nil
}

func InitializeWorker(guid uint, zookeeperAddresses []string) error {

	var err error

	if systemCoordinator, err = zookeeper.New(zookeeperAddresses); err != nil {
		return err
	}
	if dataRegistry, err = data.New(guid, "Worker", false); err != nil {
		return err
	}

	logger = NewLogger()
	properties = make(map[string]interface{})

	return nil
}

func GetSystemCoordinator() *system.Coordinator {
	return &systemCoordinator
}

func GetStorageKeyValueRegister() *storage.KeyValueRegister {
	return &keyValueRegister
}

func GetMembershipRegister() *membership.Register {
	return membershipRegister
}

func GetMembershipCoordinator() *membership.Coordinator {
	return &membershipCoordinator
}

func GetLogger() *Logger {
	return logger
}

func SetProperty(key string, value interface{}) {
	properties[key] = value
}

func GetPropertyAsString(key string) string {
	return properties[key].(string)
}

func GetPropertyAsInteger(key string) int {
	return properties[key].(int)
}

func GetDataRegistry() *data.Registry {
	return dataRegistry
}
