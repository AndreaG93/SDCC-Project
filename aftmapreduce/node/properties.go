package node

import (
	"SDCC-Project/aftmapreduce/registry"
	"SDCC-Project/cloud/amazons3"
	"SDCC-Project/cloud/zookeeper"
)

var zookeeperClient *zookeeper.Client
var logger *Logger
var properties map[string]interface{}
var amazonS3Client *amazons3.S3Client

var dataRegistry *registry.DataRegistry
var digestRegistry *registry.DigestRegistry

func Initialize(zookeeperAddresses []string) {

	zookeeperClient = zookeeper.New(zookeeperAddresses)
	amazonS3Client = amazons3.New()

	logger = NewLogger()

	dataRegistry = registry.NewDataRegistry()
	digestRegistry = registry.NewDigestRegistry()

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
