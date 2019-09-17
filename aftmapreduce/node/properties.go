package node

import (
	"SDCC-Project/aftmapreduce/cloud"
	"SDCC-Project/aftmapreduce/cloud/amazons3"
	"SDCC-Project/aftmapreduce/cloud/zookeeper"
	"SDCC-Project/aftmapreduce/registry"
	"fmt"
)

var systemCoordinator cloud.SystemCoordinator
var keyValueStorageService cloud.KeyValueStorageService
var membershipRegister cloud.MembershipRegister

var logger *Logger
var properties map[string]interface{}

var dataRegistry *registry.DataRegistry
var digestRegistry *registry.DigestRegistry

func InitializePrimary(nodeId int, zookeeperAddresses []string) error {

	var err error

	if systemCoordinator, err = zookeeper.New(zookeeperAddresses); err != nil {
		return err
	}
	if err = systemCoordinator.Initialize(); err != nil {
		return err
	}

	membershipRegister = *cloud.NewMembershipRegister()
	keyValueStorageService = amazons3.New()
	dataRegistry = registry.NewDataRegistry(fmt.Sprintf("Primary%d", nodeId), false)
	logger = NewLogger()
	properties = make(map[string]interface{})

	go membershipRegister.StartMembershipRegisterListener(systemCoordinator)

	return nil
}

func InitializeWorker(nodeId int, zookeeperAddresses []string) error {

	var err error

	if systemCoordinator, err = zookeeper.New(zookeeperAddresses); err != nil {
		return err
	}

	dataRegistry = registry.NewDataRegistry(fmt.Sprintf("Worker%d", nodeId), true)
	digestRegistry = registry.NewDigestRegistry()

	logger = NewLogger()
	properties = make(map[string]interface{})

	return nil
}

func GetSystemCoordinator() *cloud.SystemCoordinator {
	return &systemCoordinator
}

func GetKeyValueStorageService() *cloud.KeyValueStorageService {
	return &keyValueStorageService
}

func GetMembershipRegister() *cloud.MembershipRegister {
	return &membershipRegister
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

func GetDataRegistry() *registry.DataRegistry {
	return dataRegistry
}

func GetDigestRegistry() *registry.DigestRegistry {
	return digestRegistry
}
