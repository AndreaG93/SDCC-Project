package node

import (
	"SDCC-Project/aftmapreduce/cloud"
	"SDCC-Project/aftmapreduce/cloud/amazons3"
	"SDCC-Project/aftmapreduce/cloud/zookeeper"
	"SDCC-Project/aftmapreduce/registry"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
)

var systemCoordinator cloud.SystemCoordinator
var keyValueStorageService cloud.KeyValueStorageService
var membershipRegister cloud.MembershipRegister

var logger *Logger
var properties map[string]interface{}

var dataRegistry *registry.DataRegistry
var digestRegistry *registry.DigestRegistry

func InitializePrimary(zookeeperAddresses []string, nodeId int) {

	initializeNode(zookeeperAddresses)

	keyValueStorageService = amazons3.New()
	dataRegistry = registry.NewDataRegistry(fmt.Sprintf("Primary%d", nodeId), false)
}

func InitializeWorker(zookeeperAddresses []string, nodeId int) {

	initializeNode(zookeeperAddresses)

	dataRegistry = registry.NewDataRegistry(fmt.Sprintf("Worker%d", nodeId), true)
	digestRegistry = registry.NewDigestRegistry()
}

func initializeNode(zookeeperAddresses []string) {

	var err error

	systemCoordinator, err = zookeeper.New(zookeeperAddresses)
	utility.CheckError(err)

	logger = NewLogger()
	properties = make(map[string]interface{})
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
