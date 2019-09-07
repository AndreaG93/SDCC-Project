package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"fmt"
)

const (
	PendingRequestsZNodePath  = "/pending-requests"
	CompleteRequestsZNodePath = "/complete-requests"
	StatusZNodeName           = "status"
	CacheData1ZNodeName       = "data1"
	CacheData2ZNodeName       = "data2"
	InitialStatus             = "0"
	AfterMapStatus            = "1"
	AfterLocalityAwareShuffle = "2"
	AfterReduce               = "3"
	Complete                  = "4"
)

type ClientRequest struct {
	digest        string
	rootZNodePath string

	cacheData1ZNodePath string
	cacheData2ZNodePath string

	statusZNodePath          string
	completeRequestZNodePath string
}

func NewClientRequest(digest string) *ClientRequest {

	output := new(ClientRequest)

	(*output).digest = digest

	(*output).rootZNodePath = fmt.Sprintf("%s/%s", PendingRequestsZNodePath, (*output).digest)
	(*output).cacheData1ZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, CacheData1ZNodeName)
	(*output).cacheData2ZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, CacheData2ZNodeName)
	(*output).statusZNodePath = fmt.Sprintf("%s/%s", (*output).rootZNodePath, StatusZNodeName)
	(*output).completeRequestZNodePath = fmt.Sprintf("%s/%s", CompleteRequestsZNodePath, (*output).digest)

	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).rootZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).cacheData1ZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).cacheData2ZNodePath, nil, int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).statusZNodePath, []byte(InitialStatus), int32(0))
	node.GetZookeeperClient().CreateZNodeCheckingExistence((*output).completeRequestZNodePath, nil, int32(0))

	return output
}

func CheckDuplicatedClientRequest(digest string) bool {

	completeRequestZNodePath := fmt.Sprintf("%s/%s", CompleteRequestsZNodePath, digest)
	return node.GetZookeeperClient().CheckZNodeExistence(completeRequestZNodePath)
}

func (obj *ClientRequest) CheckPoint(newStatus string, data1 []byte, data2 []byte) {

	node.GetZookeeperClient().SetZNodeData((*obj).statusZNodePath, []byte(newStatus))

	if data1 != nil {
		node.GetZookeeperClient().SetZNodeData((*obj).cacheData1ZNodePath, data1)
	}
	if data2 != nil {
		node.GetZookeeperClient().SetZNodeData((*obj).cacheData2ZNodePath, data2)
	}
}

func (obj *ClientRequest) getStatus() string {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).statusZNodePath)
	return string(output)
}

func (obj *ClientRequest) GetDataFromCache1() []byte {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).cacheData1ZNodePath)
	return output
}

func (obj *ClientRequest) GetDataFromCache2() []byte {

	output, _ := node.GetZookeeperClient().GetZNodeData((*obj).cacheData2ZNodePath)
	return output
}

func (obj *ClientRequest) GetDigest() string {
	return (*obj).digest
}

func (obj *ClientRequest) GetCompleteRequestZNodePath() string {
	return (*obj).completeRequestZNodePath
}

func (obj *ClientRequest) DeletePendingRequest() {
	node.GetZookeeperClient().RemoveZNode((*obj).cacheData1ZNodePath)
	node.GetZookeeperClient().RemoveZNode((*obj).cacheData2ZNodePath)
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
