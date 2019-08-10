package node

import (
	"SDCC-Project/aftmapreduce/registry"
	"SDCC-Project/cloud/amazons3"
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/nodelogger"
)

var zookeeperClient *zookeeper.Client
var logger *nodelogger.Logger
var cache *registry.Registry
var properties map[string]interface{}
var amazonS3Client *amazons3.S3Client

func Initialize(zookeeperAddresses []string) {

	zookeeperClient = zookeeper.New(zookeeperAddresses)
	logger = nodelogger.New()
	cache = registry.New()
	properties = make(map[string]interface{})
	amazonS3Client = amazons3.New()
}

func GetZookeeperClient() *zookeeper.Client {
	return zookeeperClient
}

func GetLogger() *nodelogger.Logger {
	return logger
}

func GetCache() *registry.Registry {
	return cache
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
