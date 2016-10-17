package api

import (
	// "github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/device-manager/util"
)

var manager = bluez.NewObjectManager()
var adapters = make(map[string]*bluez.Adapter1, 0)

//GetDevices returns a list of bluetooth discovered Devices
func GetDevices() ([]Device, error) {

	objects, err := manager.GetManagedObjects()

	if err != nil {
		return nil, err
	}

	var devices = make([]Device, 0)
	for path, ifaces := range objects {
		for iface, list := range ifaces {
			switch iface {
			case "org.bluez.Device1":
				{
					deviceProperties := new(bluez.Device1Properties)
					util.MapToStruct(deviceProperties, list)
					dev := NewDevice(string(path), deviceProperties)
					devices = append(devices, *dev)
				}
			}
		}
	}

	return devices, nil
}

func getAdapter(adapterID string) (*bluez.Adapter1, error) {
	if adapters[adapterID] == nil {
		adapters[adapterID] = bluez.NewAdapter1(adapterID)
	}
	return adapters[adapterID], nil
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
	adapter, err := getAdapter(adapterID)
	if err != nil {
		return err
	}
	return adapter.StartDiscovery()
}

// StopDiscoveryOn start discovery on specified adapter
func StopDiscoveryOn(adapterID string) error {
	adapter, err := getAdapter(adapterID)
	if err != nil {
		return err
	}
	return adapter.StopDiscovery()
}
