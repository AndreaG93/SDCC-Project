package zookeeper

import (
	"SDCC-Project-WorkerNode/utility"
	"errors"
	"fmt"
	"github.com/Comcast/go-leaderelection"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

const (
	electionZNodeName         = "/election"
	electionResponseZNodeName = "/election_response"
)

func SetZNode(zookeeperConnection *zk.Conn, zNodePath string, data []byte) error {

	var actualStat *zk.Stat
	var err error

	if _, actualStat, err = zookeeperConnection.Get(zNodePath); err != nil {
		return err
	}

	if _, err := zookeeperConnection.Set(zNodePath, data, actualStat.Version); err != nil {
		return err
	}

	return nil
}

func ConnectToZookeeperServers(serverPoolAddresses []string) (*zk.Conn, <-chan zk.Event, error) {

	zkConnection, zkConnectionChannel, err := zk.Connect(serverPoolAddresses, 5*time.Second)
	if err != nil {
		return nil, nil, err
	}

	return zkConnection, zkConnectionChannel, nil
}

func checkExistenceOrGenerateZNode(zookeeperConnection *zk.Conn, zNodeName string, zNodeFlag int32) error {

	var isZNodeExisting bool
	var err error

	if zookeeperConnection == nil {
		return errors.New("'zookeeperConnection' nil")
	}

	if isZNodeExisting, _, err = zookeeperConnection.Exists(zNodeName); err != nil {
		return err
	}

	if !isZNodeExisting {
		if _, err = zookeeperConnection.Create(zNodeName, nil, zNodeFlag, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	return nil
}

func checkExistenceOrGenerateElectionZNode(zkConnection *zk.Conn) error {
	return checkExistenceOrGenerateZNode(zkConnection, electionZNodeName, zk.FlagSequence|zk.FlagEphemeral)
}

func StartLeaderElection(nodeID uint) uint {

	var err error
	var zkConnection *zk.Conn
	var zkConnectionEventChannel <-chan zk.Event
	var status leaderelection.Status
	var ok bool

	zkConnection, zkConnectionEventChannel, err = ConnectToZookeeperServers([]string{"localhost"})
	utility.CheckError(err)

	err = checkExistenceOrGenerateElectionZNode(zkConnection)
	utility.CheckError(err)

	err = checkExistenceOrGenerateZNode(zkConnection, electionResponseZNodeName, 0)
	utility.CheckError(err)

	election, err := leaderelection.NewElection(zkConnection, electionZNodeName, string(nodeID))
	utility.CheckError(err)

	zkConnectionFailedChannel := make(chan bool)

	_, _, dd, err := zkConnection.GetW(electionResponseZNodeName)
	utility.CheckError(err)

	go func() {
		for {
			evt := <-zkConnectionEventChannel

			if evt.State == zk.StateDisconnected {
				close(zkConnectionFailedChannel) // signal candidates to exit
				time.Sleep(2 * time.Second)      // give goroutines time to exit
				break
			}
		}

	}()

	go election.ElectLeader()

	for {

		select {

		case <-dd:

			fmt.Println("ho capito")
			election.Resign()

			output, _, _ := zkConnection.Get(electionResponseZNodeName)

			out, _ := strconv.ParseInt(string(output), 10, 0)

			return uint(out)

		case <-zkConnectionFailedChannel:
			fmt.Println("\t\t\tZK connection failed for candidate <", status.CandidateID, ">, exiting")
			//respCh <- ElectionResponse{false, status.CandidateID}
			election.Resign()
			return 0

		case status, ok = <-election.Status():
			if !ok {
				fmt.Println("\t\t\tChannel closed, election is terminated!!!")
				election.Resign()
				return 0
			}
			if status.Err != nil {
				fmt.Println("Received election status error <<", status.Err, ">> for candidate <", nodeID, ">.")
				election.Resign()
				return 0
			}
			fmt.Println("Candidate - ", nodeID, " - has received status message: <", status, ">.")

			if status.Role == leaderelection.Leader {
				fmt.Println("Leader is: ", nodeID)

				err := SetZNode(zkConnection, electionResponseZNodeName, []byte(string(strconv.Itoa(int(nodeID)))))
				utility.CheckError(err)

				election.EndElection() // Terminates the election and signals all followers the election is over.
				return nodeID

			}
		}
	}
}
