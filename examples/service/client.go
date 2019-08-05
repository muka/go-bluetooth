package service_example

import (
	"errors"
	"fmt"
	"strings"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

func createClient(adapterID, hwaddr, serviceID string) (err error) {

	log.Infof("Discovering devices from %s, looking for %s (serviceID:%s)", adapterID, hwaddr, serviceID)

	a, err := adapter.NewAdapter1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	log.Info("Set discovery filter")
	filter := adapter.NewDiscoveryFilter()
	filter.AddUUIDs(serviceID)

	err = a.SetDiscoveryFilter(filter.ToMap())
	if err != nil {
		return fmt.Errorf("SetDiscoveryFilter: %s", err)
	}

	log.Debug("List devices")
	devices, err := a.GetDevices()
	fail("GetDevices", err)
	for _, dev := range devices {
		err = showDeviceInfo(dev, hwaddr, serviceID)
		fail("showDeviceInfo", err)
	}

	discovery, cancel, err := api.Discover(adapterID, &filter)
	if err != nil {
		return err
	}

	defer cancel()

	for ev := range discovery {
		dev, err := device.NewDevice1(ev.Path)
		if err != nil {
			return err
		}

		err = showDeviceInfo(dev, hwaddr, serviceID)
		if err != nil {
			return err
		}

	}

	return err
}

func showDeviceInfo(dev *device.Device1, hwaddr, serviceID string) error {

	if dev == nil {
		return errors.New("Device is nil")
	}

	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("%s: Failed to get properties: %s", dev.Path(), err.Error())
	}

	// device1, err := dev.GetClient()
	// if err != nil {
	// 	return fmt.Errorf("GetClient: %s", err)
	// }

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
		log.Error(fmt.Errorf("Service UUID %s not found on %s", serviceID, hwaddr))
		return nil
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

	chars, err := dev.GetCharacteristicsList()
	if err != nil {
		return fmt.Errorf("Failed to list chars: %s", err)
	}

	descr, err := dev.GetDescriptorList()
	if err != nil {
		return fmt.Errorf("Failed to list descr: %s", err)
	}

	log.Infof("CHARS %++v", chars)
	log.Infof("DESCR %++v", descr)

	return nil
}
