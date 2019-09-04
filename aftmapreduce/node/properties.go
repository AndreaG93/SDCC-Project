package node

import (
	"SDCC-Project/aftmapreduce/cloud/amazons3"
	"SDCC-Project/aftmapreduce/cloud/zookeeper"
	"SDCC-Project/aftmapreduce/registry"
	"fmt"
)

var zookeeperClient *zookeeper.Client
var logger *Logger
var properties map[string]interface{}
var amazonS3Client *amazons3.S3Client

var dataRegistry *registry.DataRegistry
var digestRegistry *registry.DigestRegistry

func InitializePrimary(zookeeperAddresses []string, nodeId int) {

	initializeNode(zookeeperAddresses)

	amazonS3Client = amazons3.New()
	dataRegistry = registry.NewDataRegistry(fmt.Sprintf("Primary%d", nodeId), false)
}

func InitializeWorker(zookeeperAddresses []string, nodeId int) {

	initializeNode(zookeeperAddresses)

	dataRegistry = registry.NewDataRegistry(fmt.Sprintf("Worker%d", nodeId), true)
	digestRegistry = registry.NewDigestRegistry()
}

func initializeNode(zookeeperAddresses []string) {

	zookeeperClient = zookeeper.New(zookeeperAddresses)
	logger = NewLogger()

	properties = make(map[string]interface{})
}

func GetZookeeperClient() *zookeeper.Client {
	return zookeeperClient
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

func GetAmazonS3Client() *amazons3.S3Client {
	return amazonS3Client
}

func GetDataRegistry() *registry.DataRegistry {
	return dataRegistry
}

func GetDigestRegistry() *registry.DigestRegistry {
	return digestRegistry
}
