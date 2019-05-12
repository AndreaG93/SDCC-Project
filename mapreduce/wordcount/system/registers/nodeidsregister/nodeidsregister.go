package nodeidsregister

import (
	"sync"
)

var register *NodeIDsRegister
var once sync.Once

type NodeIDsRegister struct {
	ipAddresses []string

	primaries []uint
	workers   []uint
}

func GetInstance() *NodeIDsRegister {
	once.Do(func() {
		register = build()
	})
	return register
}

func build() *NodeIDsRegister {

	output := new(NodeIDsRegister)

	(*output).primaries = []uint{0, 1, 2, 3, 4}
	(*output).workers = []uint{5, 6, 7, 8, 9, 10, 11, 12, 13}
	(*output).ipAddresses = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"}

	return output
}

func (obj *NodeIDsRegister) GetNodeIDs() []uint {
	return append((*obj).primaries, (*obj).workers...)
}

func (obj *NodeIDsRegister) GetPrimaryNodeIDs() []uint {
	return (*obj).primaries
}

func (obj *NodeIDsRegister) GetWorkerNodeIDs() []uint {
	return (*obj).workers
}

func (obj *NodeIDsRegister) GetNodeIpAddress(id uint) string {
	return (*obj).ipAddresses[id]
}
