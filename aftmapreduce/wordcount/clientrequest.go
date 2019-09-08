package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"fmt"
)

const (
	PendingRequestsZNodePath  = "/pending-requests"
	CompleteRequestsZNodePath = "/complete-requests"
	StatusZNodeName           = "status"
	CacheDataZNodeName        = "data"
	InitialStatus             = "0"
	AfterMap                  = "1"
	AfterReduce               = "2"
	Complete                  = "3"
)

type ClientRequest struct {
	digest        string
	rootZNodePath string

	cacheDataZNodePath string

	statusZNodePath          string
	completeRequestZNodePath string
}

func NewClientRequest(digest string) *ClientRequest {

	output := new(ClientRequest)

	(*output).digest = digest

	(*output).rootZNodePath = fmt.Sprintf("%s/%s", PendingRequestsZNodePath, (*output).digest)
	(*output).cacheDataZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, CacheDataZNodeName)
	(*output).statusZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, StatusZNodeName)
	(*output).completeRequestZNodePath = fmt.Sprintf("%s/%s", CompleteRequestsZNodePath, (*output).digest)

	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).rootZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).cacheDataZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).statusZNodePath, []byte(InitialStatus), int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).completeRequestZNodePath, nil, int32(0))

	return output
}

func CheckDuplicatedClientRequest(digest string) bool {

	completeRequestZNodePath := fmt.Sprintf("%s/%s", CompleteRequestsZNodePath, digest)
	return node.GetZookeeperClient().CheckZNodeExistence(completeRequestZNodePath)
}

func (obj *ClientRequest) CheckPoint(newStatus string, data []byte) {

	if data != nil {
		node.GetZookeeperClient().SetZNodeData((*obj).cacheDataZNodePath, data)
	}

	node.GetZookeeperClient().SetZNodeData((*obj).statusZNodePath, []byte(newStatus))
}

func (obj *ClientRequest) getStatus() string {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).statusZNodePath)
	return string(output)
}

func (obj *ClientRequest) GetDataFromCache() []byte {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).cacheDataZNodePath)
	return output
}

func (obj *ClientRequest) GetDigest() string {
	return (*obj).digest
}

func (obj *ClientRequest) GetCompleteRequestZNodePath() string {
	return (*obj).completeRequestZNodePath
}

func (obj *ClientRequest) DeletePendingRequest() {
	node.GetZookeeperClient().RemoveZNode((*obj).cacheDataZNodePath)
	node.GetZookeeperClient().RemoveZNode((*obj).statusZNodePath)
	node.GetZookeeperClient().RemoveZNode((*obj).rootZNodePath)
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
