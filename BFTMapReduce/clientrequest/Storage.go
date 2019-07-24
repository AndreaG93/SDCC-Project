package clientrequest

import (
	"SDCC-Project/cloud/zookeeper"
	"fmt"
)

const clientsRequestsPath = "/requests"

const (
	InitialPhase     = "0"
	AfterMapPhase    = "1"
	AfterReducePhase = "2"
	Final            = "3"
)

type ClientRequest struct {
	rootZNodePath string

	dataZNodePath   string
	statusZNodePath string

	zookeeperClient *zookeeper.Client
}

func InitializationClientsRequestsPath(zookeeperClient *zookeeper.Client) {
	if !(*zookeeperClient).CheckZNodeExistence(clientsRequestsPath) {
		(*zookeeperClient).CreateZNode(clientsRequestsPath, nil, 0)
	}
}

func New(fileInputDigest string) *ClientRequest {

	output := new(ClientRequest)

	(*output).zookeeperClient = zookeeper.New([]string{"localhost:2181"})

	(*output).rootZNodePath = fmt.Sprintf("%s/%s", clientsRequestsPath, fileInputDigest)
	(*output).dataZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, "data")
	(*output).statusZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, "status")

	if !(*output).zookeeperClient.CheckZNodeExistence((*output).rootZNodePath) {
		(*output).zookeeperClient.CreateZNode((*output).rootZNodePath, nil, int32(0))
		(*output).zookeeperClient.CreateZNode((*output).statusZNodePath, []byte(InitialPhase), int32(0))
		(*output).zookeeperClient.CreateZNode((*output).dataZNodePath, nil, int32(0))
	}

	return output
}

func (obj *ClientRequest) MakeSnapshot(data []byte) {

	currentPhase, _ := (*obj).zookeeperClient.GetZNodeData((*obj).statusZNodePath)

	(*obj).zookeeperClient.SetZNodeData((*obj).dataZNodePath, data)

	switch string(currentPhase) {
	case InitialPhase:
		(*obj).zookeeperClient.SetZNodeData((*obj).statusZNodePath, []byte(AfterMapPhase))
	case AfterMapPhase:
		(*obj).zookeeperClient.SetZNodeData((*obj).statusZNodePath, []byte(AfterReducePhase))
	case AfterReducePhase:
		(*obj).zookeeperClient.SetZNodeData((*obj).statusZNodePath, []byte(Final))
	}
}

func (obj *ClientRequest) GetSnapshot(requestDigest string) (string, []byte) {

	currentPhase, _ := (*obj).zookeeperClient.GetZNodeData((*obj).statusZNodePath)
	currentData, _ := (*obj).zookeeperClient.GetZNodeData((*obj).dataZNodePath)

	return string(currentPhase), currentData
}
