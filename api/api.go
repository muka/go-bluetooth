package api

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	"github.com/muka/go-bluetooth/emitter"
)

//Exit performs a clean exit
func Exit() {
	manager, err := GetManager()
	if err == nil {
		manager.Close()
	}
	bluez.CloseConnections()
}

//GetDeviceByAddress return a Device object based on its address
func GetDeviceByAddress(adapterID, address string) (*device.Device1, error) {

	list, err := GetDeviceList(adapterID)
	if err != nil {
		return nil, err
	}

	for _, path := range list {

		dev, err := device.NewDevice1(string(path))
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
func GetDevices(adapterID string) ([]*device.Device1, error) {

	list, err := GetDeviceList(adapterID)
	if err != nil {
		return nil, err
	}

	manager, err := GetManager()
	if err != nil {
		return nil, err
	}

	objects := manager.GetObjects()

	var devices = make([]*Device, 0)
	for _, path := range list {
		object, ok := objects.Load(path)
		if !ok {
			return nil, fmt.Errorf("Path %s does not exists", path)
		}
		props := (object.(map[string]map[string]dbus.Variant))[device.Device1Interface]
		dev, err := parseDevice(path, props)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dev)
	}

	return devices, nil
}

//GetDeviceList returns a list of discovered Devices paths
func GetDeviceList(adapterID string) ([]dbus.ObjectPath, error) {

	list := []dbus.ObjectPath{}

	// in case it is left empty, return all devices avaliable?
	exists, err := AdapterExists(adapterID)
	if err != nil {
		return list, err
	}

	if !exists {
		return list, fmt.Errorf("Adapter %s not found", adapterID)
	}

	manager, err := GetManager()
	if err != nil {
		return nil, err
	}

	objects := manager.GetObjects()
	var devices []dbus.ObjectPath
	objects.Range(func(key, value interface{}) bool {
		ifaces := value.(map[string]map[string]dbus.Variant)
		path := key.(dbus.ObjectPath)
		for iface := range ifaces {
			switch iface {
			case device.Device1Interface:
				{
					devices = append(devices, path)
				}
			}
		}
		return true
	})

	return devices, nil
}

//GetGattManager return a GattManager1 instance
func GetGattManager(adapterID string) (*gatt.GattManager1, error) {

	if exists, err := AdapterExists(adapterID); !exists {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Adapter " + adapterID + " not found")
	}

	return gatt.NewGattManager1FromAdapterID(adapterID)
}

//StartDiscovery on adapter hci0
func StartCleanDiscovery() error {
	return StartCleanDiscoveryOn("hci0")
}

//StartDiscovery on adapter hci0
func StartCleanDiscoveryOn(adapterID string) error {
	err := FlushDevices(adapterID)
	if err != nil {
		return err
	}
	return StartDiscoveryOn(adapterID)
}

//StartDiscovery on adapter hci0
func StartDiscovery() error {
	return StartDiscoveryOn("hci0")
}

//StopDiscovery on adapter hci0
func StopDiscovery() error {
	return StopDiscoveryOn("hci0")
}

// StartDiscoveryOn start discovery on specified adapter
func StartDiscoveryOn(adapterID string) error {

	adapter, err := GetAdapter(adapterID)

	if err != nil {
		return err
	}

	err = adapter.StartDiscovery()

	if err != nil {
		return err
	}

	return nil
}

// StopDiscoveryOn start discovery on specified adapter
func StopDiscoveryOn(adapterID string) error {
	adapter, err := GetAdapter(adapterID)
	if err != nil {
		return err
	}
	return adapter.StopDiscovery()
}

//On add an event handler
func On(name string, fn *emitter.Callback) error {
	return emitter.On(name, fn)
}

//Off remove an event handler
func Off(name string, fn *emitter.Callback) error {
	return emitter.Off(name, fn)
}
