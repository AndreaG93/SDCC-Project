package zookeeper

import (
	"SDCC-Project/utility"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sort"
	"strings"
)

const (
	electionZNodePath = "/election"
	proposalZNodePath = "/election/"
)

type ElectionResponse struct {
	IsLeader bool
	NodeID   string
}

func (obj *Client) initializationElectionResource() {

	if !(*obj).CheckZNodeExistence(electionZNodePath) {
		(*obj).CreateZNode(electionZNodePath, nil, int32(0))
	}
}

func (obj *Client) createProposalsZNode(internetAddress string) string {

	path, err := (*obj).zooKeeperConnection.Create(proposalZNodePath, []byte(internetAddress), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	utility.CheckError(err)

	return strings.Split(path, "/")[2]
}

func (obj *Client) getCandidates() []string {

	output, _, err := (*obj).zooKeeperConnection.Children(electionZNodePath)
	utility.CheckError(err)

	sort.Strings(output)
	return output
}

func (obj *Client) GetCurrentLeaderRequestRPCInternetAddress() (string, error) {

	candidates := (*obj).getCandidates()
	if len(candidates) == 0 {
		return "", errors.New("no leader candidate")
	}

	leaderZNodePath := fmt.Sprintf("%s%s", proposalZNodePath, candidates[0])

	output, _ := (*obj).GetZNodeData(leaderZNodePath)
	return string(output), nil
}

func (obj *Client) RunAsLeaderCandidate(responseChannel chan bool, internetAddress string) {

	(*obj).initializationElectionResource()
	myProposal := (*obj).createProposalsZNode(internetAddress)
	fmt.Println("My proposal is ", myProposal)

	for {

		candidates := (*obj).getCandidates()

		if strings.EqualFold(myProposal, candidates[0]) {
			responseChannel <- true
		} else {

			zNodeToWatch := fmt.Sprintf("%s%s", proposalZNodePath, candidates[len(candidates)-2])

			fmt.Println("I'm watching ", zNodeToWatch)
			watcher := (*obj).GetZNodeWatcher(zNodeToWatch)
			<-watcher
		}
	}
}
