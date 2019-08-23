package registry

import "time"

const (
	timeout = 5 * time.Minute
)

type DataRegistry struct {
	content map[string]interface{}
}

func NewDataRegistry() *DataRegistry {

	output := new(DataRegistry)
	(*output).content = make(map[string]interface{})

	return output
}

func (obj *DataRegistry) Get(key string) interface{} {
	return (*obj).content[key]
}

func (obj *DataRegistry) Set(key string, input interface{}) {
	(*obj).content[key] = input
	go (*obj).automaticClean(key)
}

func (obj *DataRegistry) automaticClean(digest string) {

	timer1 := time.NewTimer(timeout)
	<-timer1.C

	delete((*obj).content, digest)
}
