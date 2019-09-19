package process

import (
	"SDCC-Project/aftmapreduce"
	"SDCC-Project/aftmapreduce/data"
	"SDCC-Project/aftmapreduce/membership"
	"SDCC-Project/aftmapreduce/process/property"
	"SDCC-Project/aftmapreduce/storage"
	"SDCC-Project/aftmapreduce/storage/amazons3"
	"SDCC-Project/aftmapreduce/system"
	"SDCC-Project/aftmapreduce/system/zookeeper"
	"fmt"
)

var systemCoordinator system.Coordinator
var keyValueRegister storage.KeyValueRegister
var membershipRegister *membership.Register
var membershipCoordinator membership.Coordinator
var dataRegistry *data.Registry
var properties map[string]interface{}
var logger *Logger

const (
	PrimaryProcessType = "Primary"
	WorkerProcessType  = "Worker"
)

func Initialize(processID int, processGroupID int, processType string, publicInternetAddress string, zookeeperClusterInternetAddresses []string) error {
	initializeProperties(processID, processGroupID, processType, publicInternetAddress)
	return initializeServices(zookeeperClusterInternetAddresses)
}

func initializeServices(zookeeperClusterInternetAddresses []string) error {

	var err error

	if zookeeperClient, err := zookeeper.New(zookeeperClusterInternetAddresses); err != nil {
		return err
	} else {
		systemCoordinator = zookeeperClient
		membershipCoordinator = zookeeperClient
	}

	if err = systemCoordinator.Initialize(); err != nil {
		return err
	}
	if dataRegistry, err = data.New(GetPropertyAsInteger(property.NodeID), GetPropertyAsString(property.NodeType), GetPropertyAsBoolean(property.IsDataRegisterVolatile)); err != nil {
		return err
	}
	if membershipRegister, err = membership.New(membershipCoordinator); err != nil {
		return err
	}
	if GetPropertyAsBoolean(property.CanAccessToDFS) {
		keyValueRegister = amazons3.New()
	}

	logger = NewLogger()

	return nil
}

func initializeProperties(processID int, processGroupID int, processType string, publicInternetAddress string) {

	properties = make(map[string]interface{})

	SetProperty(property.NodeID, processID)
	SetProperty(property.NodeType, processType)
	SetProperty(property.InternetAddress, publicInternetAddress)

	switch processType {
	case PrimaryProcessType:

		SetProperty(property.WordCountRequestRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountRequestRPCBasePort+processID))
		SetProperty(property.IsDataRegisterVolatile, true)
		SetProperty(property.CanAccessToDFS, true)

	case WorkerProcessType:

		SetProperty(property.NodeGroupID, processGroupID)
		SetProperty(property.WordCountMapRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountMapTaskRPCBasePort+processID))
		SetProperty(property.WordCountReduceRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountReduceTaskRPCBasePort+processID))
		SetProperty(property.WordCountReceiveRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountReceiveRPCBasePort+processID))
		SetProperty(property.WordCountSendRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountSendRPCBasePort+processID))
		SetProperty(property.WordCountRetrieveRPCFullAddress, fmt.Sprintf("%s:%d", "localhost", aftmapreduce.WordCountRetrieverRPCBasePort+processID))

		SetProperty(property.IsDataRegisterVolatile, false)
		SetProperty(property.CanAccessToDFS, false)
	}
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

func GetPropertyAsBoolean(key string) bool {
	return properties[key].(bool)
}

func GetDataRegistry() *data.Registry {
	return dataRegistry
}