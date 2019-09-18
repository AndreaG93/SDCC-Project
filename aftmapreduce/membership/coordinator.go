package membership

type Coordinator interface {
	WaitUntilProcessMembershipChanges() error
	GetProcessMembershipTable() (map[int]map[int]string, error)
	RegisterNewWorkerProcess(processId int, processGroupId int, processPublicInternetAddress string) error
}
