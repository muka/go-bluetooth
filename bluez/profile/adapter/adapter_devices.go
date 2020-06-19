package adapter

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/util"
)

//GetDeviceByAddress return a Device object based on its address
func (a *Adapter1) GetDeviceByAddress(address string) (*device.Device1, error) {

	list, err := a.GetDeviceList()
	if err != nil {
		return nil, err
	}

	for _, path := range list {

		dev, err := device.NewDevice1(path)
		if err != nil {
			return nil, err
		}

		if dev.Properties.Address == address {
			return dev, nil
		}
	}

	return nil, nil
}

//GetDevices returns a list of bluetooth discovered Devices
func (a *Adapter1) GetDevices() ([]*device.Device1, error) {

	list, err := a.GetDeviceList()
	if err != nil {
		return nil, err
	}

	om, err := bluez.GetObjectManager()
	if err != nil {
		return nil, err
	}

	objects, err := om.GetManagedObjects()
	if err != nil {
		return nil, err
	}

	devices := []*device.Device1{}

	for _, path := range list {

		object, ok := objects[path]
		if !ok {
			return nil, fmt.Errorf("Path %s does not exists", path)
		}

		props := object[device.Device1Interface]
		dev, err := parseDevice(path, props)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dev)
	}

	// log.Debugf("%d cached devices", len(devices))
	return devices, nil
}

// GetDeviceList returns a list of cached device paths
func (a *Adapter1) GetDeviceList() ([]dbus.ObjectPath, error) {

	om, err := bluez.GetObjectManager()
	if err != nil {
		return nil, err
	}

	objects, err := om.GetManagedObjects()
	if err != nil {
		return nil, err
	}

	devices := []dbus.ObjectPath{}
	for path, ifaces := range objects {
		for iface := range ifaces {
			switch iface {
			case device.Device1Interface:
				{
					if strings.Contains(string(path), string(a.Path())) {
						devices = append(devices, path)
					}
				}
			}
		}
	}

	return devices, nil
}

// FlushDevices removes device from bluez cache
func (a *Adapter1) FlushDevices() error {

	devices, err := a.GetDevices()
	if err != nil {
		return err
	}

	for _, dev := range devices {
		err = a.RemoveDevice(dev.Path())
		if err != nil {
			return fmt.Errorf("FlushDevices.RemoveDevice %s: %s", dev.Path(), err)
		}
	}

	return nil
}

// ParseDevice parse a Device from a ObjectManager map
func parseDevice(path dbus.ObjectPath, propsMap map[string]dbus.Variant) (*device.Device1, error) {

	dev, err := device.NewDevice1(path)
	if err != nil {
		return nil, err
	}

	err = util.MapToStruct(dev.Properties, propsMap)
	if err != nil {
		return nil, err
	}

	return dev, nil
}
