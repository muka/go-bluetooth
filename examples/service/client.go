package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/linux"
)

func createClient(adapterID, name, path string) error {

	log.Info("Discovering devices")

	btmgmt := linux.NewBtMgmt(adapterID)

	// turn off/on
	err := btmgmt.Reset()
	if err != nil {
		return err
	}

	adapter := profile.NewAdapter1(clientAdapterID)
	err = adapter.StartDiscovery()
	if err != nil {
		log.Errorf("Failed to start discovery: %s", err.Error())
		return err
	}

	devices, err := api.GetDevices()
	if err != nil {
		return err
	}

	for _, d := range devices {
		err := adapter.RemoveDevice(d.Path)
		if err != nil {
			log.Warnf("Cannot remove %s : %s", d.Path, err.Error())
		}
	}

	log.Infof("Start discovery..")
	api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {

		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		if discoveryEvent.Status == api.DeviceAdded {
			showDeviceInfo(discoveryEvent.Device)
		}

	}))

	return nil
}

func showDeviceInfo(dev *api.Device) {
	if dev == nil {
		return
	}
	props, err := dev.GetProperties()
	if err != nil {
		log.Errorf("%s: Failed to get properties: %s", dev.Path, err.Error())
		return
	}
	log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
}
