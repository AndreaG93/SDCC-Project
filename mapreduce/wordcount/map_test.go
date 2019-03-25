package wordcount

import "testing"

func Test_MapService(t *testing.T) {

	var err error

	mapInput := MapInput{"map_test", 5}
	mapOutput := MapOutput{}
	mapObject := Map{}

	if err = mapObject.Execute(mapInput, &mapOutput); err != nil {
		panic(err)
	}
}
