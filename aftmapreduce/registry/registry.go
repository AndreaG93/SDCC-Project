package registry

type Registry struct {
	content map[string]interface{}
}

func New() *Registry {

	output := new(Registry)
	(*output).content = make(map[string]interface{})

	return output
}

func (obj *Registry) Get(key string) interface{} {
	return (*obj).content[key]
}

func (obj *Registry) Set(key string, input interface{}) {
	(*obj).content[key] = input
}
