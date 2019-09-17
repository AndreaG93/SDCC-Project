package cloud

type SystemCoordinator interface {
	WaitUntilProcessMembershipChanges() error
	UpdateProcessMembershipRegister() (map[int]map[int]string, error)
	RegisterNewWorkerProcess(processId int, processGroupId int, processPublicInternetAddress string) error
	RegisterClientRequest(guid string, initialStatus uint8) error
	DeletePendingRequest(guid string) error
	GetAllPendingClientRequestGuid() ([]string, error)
	ClientRequestExist(guid string) (bool, error)
	GetClientRequestInformation(guid string) (uint8, []byte, error)
	UpdateClientRequestStatusBackup(guid string, status uint8, data []byte) error
	WaitUntilLeader(myOwnPublicInternetAddress string) error
	RegisterClientRequestAsComplete(guid string, outputGuid string) error
	Initialize() error
}
