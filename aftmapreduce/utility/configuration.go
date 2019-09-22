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

func (obj *NodeConfiguration) Load(path string) error {

	if configurationFile, err := os.Open(path); err != nil {
		return err
	} else {
		defer func() {
			CheckError(configurationFile.Close())
		}()

		decoder := json.NewDecoder(configurationFile)
		return decoder.Decode(obj)
	}
}
