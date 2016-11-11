package api

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/op/go-logging"
	"github.com/tj/go-debug"

	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/bluez/profile"
	"github.com/muka/bluez-client/emitter"
)

var logger = logging.MustGetLogger("api")

var dbg = debug.Debug("bluez:api")

//Exit performs a clean exit
func Exit() {
	GetManager().Unregister()
	GetManager().Close()
	dbg("Bye.")
}

//GetDevices returns a list of bluetooth discovered Devices
func GetDevices() ([]Device, error) {

	objects, err := GetManager().GetManagedObjects()

	if err != nil {
		return nil, err
	}

	var devices = make([]Device, 0)
	for path, ifaces := range objects {
		for iface, props := range ifaces {
			switch iface {
			case bluez.Device1Interface:
				{
					dev, err := ParseDevice(path, props)
					if err != nil {
						return nil, err
					}
					devices = append(devices, *dev)
				}
			}
		}
	}

	return devices, nil
}

//AdapterExists checks if an adapter is available
func AdapterExists(adapterID string) (bool, error) {
	objects, err := GetManager().GetManagedObjects()
	if err != nil {
		return false, err
	}
	path := dbus.ObjectPath("/org/bluez/" + adapterID)
	return objects[path] != nil, nil
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

	// register for manager signals, return chan *dbus.Signal
	err = WatchManagerChanges()
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
func On(name string, fn Callback) {
	emitter.On(name, func(ev emitter.Event) {
		fn(ev)
	})
}

//Off remove an event handler
func Off(name string) {
	emitter.Off(name)
}
