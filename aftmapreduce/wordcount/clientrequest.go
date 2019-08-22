package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"fmt"
)

const (
	PendingRequestsZNodePath  = "/pending-requests"
	CompleteRequestsZNodePath = "/complete-requests"
	StatusZNodeName           = "status"
	DataZNodeName             = "data"
	InitialStatus             = "0"
	AfterMapStatus            = "1"
	AfterLocalityAwareShuffle = "2"
	AfterReduce               = "3"
	Complete                  = "4"
)

type ClientRequest struct {
	digest                   string
	rootZNodePath            string
	dataZNodePath            string
	statusZNodePath          string
	completeRequestZNodePath string
}

func NewClientRequest(digest string) *ClientRequest {

	output := new(ClientRequest)

	(*output).digest = digest

	(*output).rootZNodePath = fmt.Sprintf("%s/%s", PendingRequestsZNodePath, (*output).digest)
	(*output).dataZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, DataZNodeName)
	(*output).statusZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, StatusZNodeName)
	(*output).completeRequestZNodePath = fmt.Sprintf("%s/%s", CompleteRequestsZNodePath, (*output).digest)

	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).rootZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).dataZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).statusZNodePath, []byte(InitialStatus), int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).completeRequestZNodePath, nil, int32(0))

	return output
}

func (obj *ClientRequest) CheckPoint(newStatus string, data []byte) {

	node.GetZookeeperClient().SetZNodeData((*obj).statusZNodePath, []byte(newStatus))
	node.GetZookeeperClient().SetZNodeData((*obj).dataZNodePath, data)
}

func (obj *ClientRequest) getStatus() string {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).statusZNodePath)
	return string(output)
}

func (obj *ClientRequest) GetData() []byte {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).dataZNodePath)
	return output
}

func (obj *ClientRequest) GetDigest() string {
	return (*obj).digest
}

func (obj *ClientRequest) GetCompleteRequestZNodePath() string {
	return (*obj).completeRequestZNodePath
}

func InitNeededZNodePathsToManageClientRequests() {

	node.GetZookeeperClient().CreateZNodeCheckingExistence(PendingRequestsZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence(CompleteRequestsZNodePath, nil, int32(0))
}

func GetPendingClientsRequests() []*ClientRequest {

	zookeeperClient := node.GetZookeeperClient()

	pendingClientRequests := zookeeperClient.GetChildrenList(PendingRequestsZNodePath)
	output := make([]*ClientRequest, len(pendingClientRequests))

	for index, clientRequestName := range pendingClientRequests {
		output[index] = NewClientRequest(clientRequestName)
	}

	return output
}
