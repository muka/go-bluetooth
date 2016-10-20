package main

import (
	"github.com/muka/bluez-client/api"
	"github.com/muka/bluez-client/emitter"
	"github.com/muka/bluez-client/util"
)

var logger = util.NewLogger("main")

func main() {

	defer api.Exit()

	adapterID := "hci0"

	if exists, err := api.AdapterExists(adapterID); !exists {
		if err != nil {
			panic(err)
		}
		logger.Println("Waiting for adapter hci0")
		emitter.On("adapter", func(ev emitter.Event) {
			info := ev.GetData().(api.AdapterEvent)

			if info.Status == api.DeviceAdded {
				logger.Printf("Adapter %s added\n", info.Name)
				discoverDevices(info.Name)

			} else {
				logger.Printf("Adapter %s removed\n", info.Name)
			}
		})
	} else {
		discoverDevices(adapterID)
	}

	select {}
}

func discoverDevices(adapterID string) {

	logger.Printf("Starting discovery on adapter %s\n", adapterID)

	var err error
	devices, err := api.GetDevices()
	if err != nil {
		panic(err)
	}

	hci0, err := api.GetAdapter(adapterID)
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		logger.Printf("Dropping %s", device.Path)
		go hci0.RemoveDevice(device.Path)
	}

	err = api.StopDiscovery()
	if err != nil {
		logger.Println(err)
	}

	err = api.StartDiscovery()
	if err != nil {
		panic(err)
	}

	emitter.On("discovery", func(ev emitter.Event) {
		info := ev.GetData().(api.DiscoveredDeviceEvent)
		if info.Status == api.DeviceAdded {

			name := info.Device.GetProperties().Name

			if name != "MI Band 2" {
				return
			}

			logger.Printf("Found device %s, connecting profiles", name)

			connectProfiles(info.Device)

			// 		logger.Printf("Found device %s, watching for property change", info.Device.GetProperties().Name)
			//
			// 		info.Device.On("change", func(ev emitter.Event) {
			// 			changed := ev.GetData().(api.PropertyChanged)
			// 			logger.Printf("%s: set %s = %s", info.Device.GetProperties().Name, changed.Field, changed.Value)
			// 		})
			//
		} else {
			logger.Printf("Removed device %s", info.Path)
		}
	})

}

func connectProfiles(dev *api.Device) {

	for _, uuid := range dev.GetProperties().UUIDs {

		logger.Printf("Connecting profile %s", uuid)

	}
}
