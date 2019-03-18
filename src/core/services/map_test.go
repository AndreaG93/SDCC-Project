package services

import (
	"os"
	"testing"
)

func Test_MapService(t *testing.T) {

	mapInput := MapInput{"Andrea Graziani supera la prova progettuale di SDCC.", 5}
	mapOutput := MapOutput{}
	mapObject := Map{}

	if myError := mapObject.Execute(mapInput, &mapOutput); myError != nil {
		os.Exit(1)
	}
}
