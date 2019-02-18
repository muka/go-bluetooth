package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	log "github.com/sirupsen/logrus"
)

func createClient(adapterID, hwaddr, serviceID string) (err error) {

	log.Infof("Discovering devices from %s, looking for %s (serviceID:%s)", adapterID, hwaddr, serviceID)

	adapter := profile.NewAdapter1(adapterID)

	err = adapter.StartDiscovery()
	if err != nil {
		log.Errorf("Failed to start discovery: %s", err.Error())
		return err
	}

	devices, err := api.GetDevices()
	fail("GetDevices", err)
	for _, dev := range devices {
		err = showDeviceInfo(&dev, hwaddr, serviceID)
		fail("showDeviceInfo", err)
	}

	err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		if discoveryEvent.Status == api.DeviceAdded {
			err = showDeviceInfo(discoveryEvent.Device, hwaddr, serviceID)
			fail("showDeviceInfo", err)
		}
	}))

	return err
}

func showDeviceInfo(dev *api.Device, hwaddr, serviceID string) error {

	if dev == nil {
		return errors.New("Device is nil")
	}

	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("%s: Failed to get properties: %s", dev.Path, err.Error())
	}

	if strings.ToLower(hwaddr) != strings.ToLower(props.Address) {
		// log.Debugf("Skip device name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
		return nil
	}

	serviceID = strings.ToLower(serviceID)

	log.Infof("Found device name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)

	var found bool
	for _, uuid := range props.UUIDs {
		if strings.ToLower(uuid) == serviceID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Service UUID %s not found on %s", serviceID, hwaddr)
	}

	if !props.Paired {
		log.Debugf("Pairing to %s...", props.Name)
		err = dev.Pair()
		if err != nil {
			return fmt.Errorf("Pair: %s", err)
		}
	}

	log.Debugf("Connecting to %s...", props.Name)
	err = dev.Connect()
	if err != nil {
		return fmt.Errorf("Connect: %s", err)
	}

	log.Infof("Found UUID %s", serviceID)

	chars, err := dev.GetCharsList()
	if err != nil {
		return fmt.Errorf("Service UUID %s not found on %s", serviceID, hwaddr)
	}

	log.Infof("CHARS %++v", chars)

	return nil
}
