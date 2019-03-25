package main

import (
	"core/services"
	"core/utility"
)

func main() {

	go func() {
		if err := utility.StartRemoteService(services.Map{}, 2000); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := utility.StartRemoteService(services.Reduce{}, 3000); err != nil {
			panic(err)
		}
	}()
}
