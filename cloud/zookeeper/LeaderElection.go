package zookeeper

import (
	"SDCC-Project/utility"
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
		(*obj).CreateZNode(electionZNodePath, int32(0))
	}
}

func (obj *Client) createProposalsZNode() string {

	path, err := (*obj).zooKeeperConnection.Create(proposalZNodePath, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	utility.CheckError(err)

	return strings.Split(path, "/")[2]
}

func (obj *Client) getCandidates() []string {

	output, _, err := (*obj).zooKeeperConnection.Children(electionZNodePath)
	utility.CheckError(err)

	sort.Strings(output)
	return output
}

func (obj *Client) RunAsLeaderCandidate(responseChannel chan bool) {

	(*obj).initializationElectionResource()
	myProposal := (*obj).createProposalsZNode()
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
