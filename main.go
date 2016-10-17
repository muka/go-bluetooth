package main

import (
	"github.com/muka/bluez-client/api"
	"github.com/muka/bluez-client/util"
)

var log = util.NewLogger("main")

func main() {
	//
	// err := api.StartDiscovery()
	// if err != nil {
	// 	panic(err)
	// }

	api.On("device", func(ev api.Event) {
		log.Printf("Got event %s: %s", ev.Name, ev.Data)
	})

	api.Emit("device", "Test")

}
