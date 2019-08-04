package api

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/util"
)

// FlushDevices clears removes device and clears it from cached devices
func FlushDevices(adapterID string) error {

	adapter, err := GetAdapter(adapterID)
	if err != nil {
		return err
	}

	devices, err := GetDevices(adapterID)
	if err != nil {
		return err
	}
	for _, dev := range devices {
		err := ClearDevice(dev)
		if err != nil {
			return fmt.Errorf("FlushDevices.ClearDevice %s: %s", dev.Path, err)
		}
		err = adapter.RemoveDevice(dbus.ObjectPath(dev.Path))
		if err != nil {
			return fmt.Errorf("FlushDevices.RemoveDevice %s: %s", dev.Path, err)
		}
	}
	return nil
}

// ParseDevice parse a Device from a ObjectManager map
func parseDevice(path dbus.ObjectPath, propsMap map[string]dbus.Variant) (*device.Device1, error) {

	dev, err := device.NewDevice1(string(path))
	if err != nil {
		return nil, err
	}

	util.MapToStruct(dev.Properties, propsMap)

	return dev, nil
}
