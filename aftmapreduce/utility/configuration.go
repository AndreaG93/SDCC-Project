package utility

import (
	"encoding/json"
	"os"
)

type NodeConfiguration struct {
	ZookeeperServersPrivateIPs []string
	NodeID                     int
	NodeGroupID                int
	NodeClass                  string
}

func (obj *NodeConfiguration) Load(path string) {

	configurationFile, err := os.Open(path)
	CheckError(err)

	defer func() { CheckError(configurationFile.Close()) }()

	decoder := json.NewDecoder(configurationFile)
	CheckError(decoder.Decode(obj))
}
