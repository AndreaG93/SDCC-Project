package zookeeper

import (
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
)

type clientRequestZNodeData struct {
	Status uint8
	Data   []byte
}

func (obj *Client) GetAllPendingClientRequestGuid() ([]string, error) {

	pendingClientRequests, _, err := (*obj).zooKeeperConnection.Children(pendingClientRequestsZNodePath)
	return pendingClientRequests, err
}

func (obj *Client) DeletePendingRequest(guid string) error {

	zNodePath := fmt.Sprintf("%s/%s", pendingClientRequestsZNodePath, guid)
	return (*obj).removeZNodeCheckingExistence(zNodePath)
}

func (obj *Client) ClientRequestExist(guid string) (bool, error) {

	zNodePath := fmt.Sprintf("%s/%s", completeClientRequestsZNodePath, guid)

	isZNodeExistent, _, err := (*obj).zooKeeperConnection.Exists(zNodePath)
	return isZNodeExistent, err
}

func (obj *Client) RegisterClientRequest(guid string, initialStatus uint8) error {

	var err error

	zNodeData := new(clientRequestZNodeData)
	(*zNodeData).Data = make([]byte, 0)
	(*zNodeData).Status = initialStatus

	zNodePath := fmt.Sprintf("%s/%s", pendingClientRequestsZNodePath, guid)
	finalizedRequestZNodePath := fmt.Sprintf("%s/%s", completeClientRequestsZNodePath, guid)

	if err = (*obj).createZNodeCheckingExistence(zNodePath, utility.Encode(zNodeData), int32(0)); err != nil {
		return err
	}
	if err = (*obj).createZNodeCheckingExistence(finalizedRequestZNodePath, nil, int32(0)); err != nil {
		return err
	}

	return nil
}

func (obj *Client) UpdateClientRequestStatusBackup(guid string, status uint8, data []byte) error {

	zNodePath := fmt.Sprintf("%s/%s", pendingClientRequestsZNodePath, guid)

	zNodeData := new(clientRequestZNodeData)
	(*zNodeData).Data = data
	(*zNodeData).Status = status

	return (*obj).setZNodeData(zNodePath, utility.Encode(zNodeData))
}

func (obj *Client) GetClientRequestInformation(guid string) (uint8, []byte, error) {

	zNodePath := fmt.Sprintf("%s/%s", pendingClientRequestsZNodePath, guid)

	clientRequestData := clientRequestZNodeData{}

	if rawData, _, err := (*obj).getZNodeData(zNodePath); err != nil {
		return 0, nil, err
	} else {
		utility.Decode(rawData, &clientRequestData)
		return clientRequestData.Status, clientRequestData.Data, nil
	}
}

func (obj *Client) RegisterClientRequestAsComplete(guid string, outputGuid string) error {

	zNodePath := fmt.Sprintf("%s/%s", completeClientRequestsZNodePath, guid)
	return (*obj).setZNodeData(zNodePath, []byte(outputGuid))
}

func (obj *Client) WaitForClientRequestCompletion(guid string) (string, error) {

	zNodePath := fmt.Sprintf("%s/%s", completeClientRequestsZNodePath, guid)

	for {

		if data, watcher, err := (*obj).getZNodeData(zNodePath); err == nil {
			if data == nil {
				<-watcher
				continue
			} else {
				return string(data), nil
			}
		} else {
			return "", err
		}
	}
}
