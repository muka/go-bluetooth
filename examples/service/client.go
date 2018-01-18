package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
)

func createClient(name string, path string) {

	log.Info("Discovering devices")

	adapter := profile.NewAdapter1(clientAdapterID)

	err := adapter.StartDiscovery()
	if err != nil {
		log.Errorf("Failed to start discovery: %s", err.Error())
	}

	api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {

		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device

		if dev == nil {
			log.Infof("Device removed %s", dev.Path)
			return
		}

		log.Infof("Device found %s", dev.Path)

	}))

}
