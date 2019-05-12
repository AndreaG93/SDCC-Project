package nodestatusregister

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/registers/nodeidsregister"
	"sync"
)

var register *NodeStatusRegister
var once sync.Once

type NodeStatusRegister struct {
	value map[uint]bool
}

func GetInstance() *NodeStatusRegister {
	once.Do(func() {
		register = build()
	})
	return register
}

func build() *NodeStatusRegister {

	nodeIDsRegister := nodeidsregister.GetInstance()

	output := new(NodeStatusRegister)

	(*output).value = make(map[uint]bool)

	for id := range nodeIDsRegister.GetNodeIDs() {
		(*output).value[uint(id)] = false
	}

	return output
}

func (obj *NodeStatusRegister) SetNodeStatusAsOnline(id uint) {
	(*obj).value[id] = true
}

func (obj *NodeStatusRegister) SetNodeStatusAsOffline(id uint) {
	(*obj).value[id] = false
}

func (obj *NodeStatusRegister) IsNodeOnline(id uint) bool {
	return (*obj).value[id] == true
}

func (obj *NodeStatusRegister) IsNodeOffline(id uint) bool {
	return (*obj).value[id] == false
}
