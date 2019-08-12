package registry

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
}
