package api

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus"

	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
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
func GetDeviceByAddress(address string) (*Device, error) {
	list, err := GetDeviceList()
	if err != nil {
		return nil, err
	}
	for _, path := range list {

		dev, err := NewDevice(string(path))
		if err != nil {
			return nil, err
		}

		dev.lock.RLock()
		// get current Properties pointer (can be changed by other goroutine)
		props := dev.Properties
		dev.lock.RUnlock()

		props.Lock()
		prop_address := props.Address
		props.Unlock()
		if prop_address == address {
			return dev, nil
		}
	}
	return nil, nil
}

//GetDevices returns a list of bluetooth discovered Devices
func GetDevices() ([]*Device, error) {

	list, err := GetDeviceList()
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
			return nil, errors.New("Path " + string(path) + " does not exists.")
		}
		props := (object.(map[string]map[string]dbus.Variant))[device.Device1Interface]
		dev, err := ParseDevice(path, props)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dev)
	}

	return devices, nil
}

//GetDeviceList returns a list of discovered Devices paths
func GetDeviceList() ([]dbus.ObjectPath, error) {

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

//AdapterExists checks if an adapter is available
func AdapterExists(adapterID string) (bool, error) {

	manager, err := GetManager()
	if err != nil {
		return false, err
	}

	objects := manager.GetObjects()

	path := dbus.ObjectPath("/org/bluez/" + adapterID)
	_, exists := objects.Load(path)

	return exists, nil
}

//GetAdapter return an adapter object instance
func GetAdapter(adapterID string) (*adapter.Adapter1, error) {

	if exists, err := AdapterExists(adapterID); !exists {
		if err != nil {
			return nil, fmt.Errorf("AdapterExists: %s", err)
		}
		return nil, errors.New("Adapter " + adapterID + " not found")
	}

	return adapter.NewAdapter1FromAdapterID(adapterID)
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
