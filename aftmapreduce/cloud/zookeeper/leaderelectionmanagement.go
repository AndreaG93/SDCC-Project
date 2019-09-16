package zookeeper

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sort"
	"strings"
)

type ElectionResponse struct {
	IsLeader bool
	NodeID   string
}

func (obj *Client) WaitUntilLeader(myOwnPublicInternetAddress string) error {

	var err error

	if myProposal, err := (*obj).createProposalsZNode(myOwnPublicInternetAddress); err == nil {

		for {

			if candidates, err := (*obj).getCandidates(); err != nil {
				break
			} else {
				if strings.EqualFold(myProposal, candidates[0]) {
					return nil
				} else {

					zNodeToWatch := fmt.Sprintf("%s/%s", electionZNodePath, candidates[len(candidates)-2])

					if _, _, watcher, err := (*obj).zooKeeperConnection.GetW(zNodeToWatch); err != nil {
						break
					} else {
						<-watcher
					}
				}
			}
		}
	}

	return err
}

func (obj *Client) createProposalsZNode(internetAddress string) (string, error) {

	if path, err := (*obj).zooKeeperConnection.Create(electionZNodePath, []byte(internetAddress), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll)); err == nil {
		return strings.Split(path, "/")[2], nil
	} else {
		return "", nil
	}
}

func (obj *Client) getCandidates() ([]string, error) {

	if output, _, err := (*obj).zooKeeperConnection.Children(electionZNodePath); err == nil {
		sort.Strings(output)
		return output, nil
	} else {
		return nil, err
	}
}

func (obj *Client) GetCurrentLeaderPublicInternetAddress() (string, error) {

	var err error

	if candidates, err := (*obj).getCandidates(); err == nil {

		if len(candidates) == 0 {
			return "", errors.New("no leader candidate")
		} else {
			leaderZNodePath := fmt.Sprintf("%s%s", electionZNodePath, candidates[0])
			if output, _, err := (*obj).zooKeeperConnection.Get(leaderZNodePath); err == nil {
				return string(output), nil
			}
		}
	}

	return "", err
}
