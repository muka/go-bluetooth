package api

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/op/go-logging"
	"github.com/tj/go-debug"

	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
)

var logger = logging.MustGetLogger("api")

var dbg = debug.Debug("bluez:api")

//Exit performs a clean exit
func Exit() {
	GetManager().Close()
	dbg("Bye.")
}

//GetDeviceByAddress return a Device object based on its address
func GetDeviceByAddress(address string) (*Device, error) {
	list, err := GetDeviceList()
	if err != nil {
		return nil, err
	}
	for _, path := range list {
		// dbg("Check device %s", path)
		dev := NewDevice(string(path))
		if dev.Properties.Address == address {
			// dbg("Address found")
			return dev, nil
		}
	}
	return nil, nil
}

//GetDevices returns a list of bluetooth discovered Devices
func GetDevices() ([]Device, error) {

	list, err := GetDeviceList()
	if err != nil {
		return nil, err
	}

	objects := GetManager().GetObjects()

	var devices = make([]Device, 0)
	for _, path := range list {
		props := (*objects)[path][bluez.Device1Interface]
		dev, err := ParseDevice(path, props)
		if err != nil {
			return nil, err
		}
		devices = append(devices, *dev)
	}

	return devices, nil
}

//GetDeviceList returns a list of discovered Devices paths
func GetDeviceList() ([]dbus.ObjectPath, error) {

	objects := GetManager().GetObjects()
	var devices []dbus.ObjectPath
	for path, ifaces := range *objects {
		for iface := range ifaces {
			switch iface {
			case bluez.Device1Interface:
				{
					devices = append(devices, path)
				}
			}
		}
	}

	return devices, nil
}

//AdapterExists checks if an adapter is available
func AdapterExists(adapterID string) (bool, error) {

	objects := GetManager().GetObjects()

	path := dbus.ObjectPath("/org/bluez/" + adapterID)
	_, exists := (*objects)[path]

	dbg("Adapter %s exists ? %t", adapterID, exists)
	return exists, nil
}

//GetAdapter return an adapter object instance
func GetAdapter(adapterID string) (*profile.Adapter1, error) {

	if exists, err := AdapterExists(adapterID); !exists {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Adapter " + adapterID + " not found")
	}

	return profile.NewAdapter1(adapterID), nil
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
func On(name string, fn *emitter.Callback) {
	emitter.On(name, fn)
}

//Off remove an event handler
func Off(name string, fn *emitter.Callback) {
	emitter.Off(name, fn)
}
