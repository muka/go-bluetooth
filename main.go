package main

import (
	"github.com/muka/bluez-client/api"
	"github.com/muka/bluez-client/emitter"
	"github.com/muka/bluez-client/util"
)

var log = util.NewLogger("main")

func main() {

	defer api.Exit()

	var err error

	devices, err := api.GetDevices()
	if err != nil {
		panic(err)
	}

	hci0, err := api.GetAdapter("hci0")
	if err != nil {
		panic(err)
	}

	for _, device := range devices {
		log.Printf("Dropping %s", device.Path)
		go hci0.RemoveDevice(device.Path)
	}

	err = api.StopDiscovery()
	if err != nil {
		log.Println(err)
	}

	err = api.StartDiscovery()
	if err != nil {
		panic(err)
	}

	emitter.On("discovery", func(ev emitter.Event) {
		info := ev.GetData().(api.DiscoveredDevice)
		if info.Status == api.DeviceAdded {
			log.Printf("Found device %s", info.Device.GetProperties().Name)
		} else {
			log.Printf("Removed device %s", info.Path)
		}

	})

	log.Println("Waiting...")
	select {}
}
