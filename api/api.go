package api

import (
	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/emitter"
	utilb "github.com/muka/bluez-client/util"
)

var log = utilb.NewLogger("api")

var manager = bluez.NewObjectManager()
var adapters = make(map[string]*bluez.Adapter1, 0)

var registrations map[string]*dbus.Signal

//Exit performs a clean exit
func Exit() {
	GetManager().Unregister()
	GetManager().Close()
	for _, adapter := range adapters {
		// adapter.Unregister()
		adapter.Close()
	}
	log.Println("Bye.")
}

//GetManager return the object manager reference
func GetManager() *bluez.ObjectManager {
	return manager
}

//GetDevices returns a list of bluetooth discovered Devices
func GetDevices() ([]Device, error) {

	objects, err := manager.GetManagedObjects()

	if err != nil {
		return nil, err
	}

	var devices = make([]Device, 0)
	for path, ifaces := range objects {
		for iface, props := range ifaces {
			switch iface {
			case bluez.Device1Interface:
				{
					dev := ParseDevice(path, props)
					devices = append(devices, *dev)
				}
			}
		}
	}

	return devices, nil
}

//GetAdapter return an adapter object instance
func GetAdapter(adapterID string) (*bluez.Adapter1, error) {
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

	adapter, err := GetAdapter(adapterID)

	if err != nil {
		return err
	}

	err = adapter.StartDiscovery()

	if err != nil {
		return err
	}

	// register for manager signals, return chan *dbus.Signal
	err = WatchDiscovery()
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

// UnwatchDiscovery regitster for signals from the ObjectManager
func UnwatchDiscovery() error {
	return GetManager().Unregister()
}

// DeviceStatus indicate the status of a device
type DeviceStatus int

const (
	//DeviceAdded indicates the device interface appeared
	DeviceAdded DeviceStatus = iota
	//DeviceRemoved indicates the device interface disappeared
	DeviceRemoved
)

//DiscoveredDevice contains detail regarding an ongoing discovery operation
type DiscoveredDevice struct {
	Path   string
	Status DeviceStatus
	Device *Device
}

// WatchDiscovery regitster for signals from the ObjectManager
func WatchDiscovery() error {

	channel, err := GetManager().Register()
	if err != nil {
		return err
	}

	log.Println("Waiting for devices discovery")
	go (func() {
		for {

			if channel == nil {
				log.Println("Quitting discovery listener")
				break
			}

			v := <-channel

			switch v.Name {
			case bluez.InterfacesRemoved:

				log.Printf("Removed %s %s", v.Name, v.Path)

				path := v.Body[0].(dbus.ObjectPath)
				ifaces := v.Body[1].([]string)
				for _, iF := range ifaces {
					if iF == bluez.Device1Interface {
						log.Printf("%s : %s", path, ifaces)
						devInfo := DiscoveredDevice{string(path), DeviceRemoved, nil}
						emitter.Emit("discovery", devInfo)
					}
				}

				break
			case bluez.InterfacesAdded:

				log.Printf("Added %s %s", v.Name, v.Path)

				path := v.Body[0].(dbus.ObjectPath)
				props := v.Body[1].(map[string]map[string]dbus.Variant)
				dev := ParseDevice(path, props[bluez.Device1Interface])
				devInfo := DiscoveredDevice{string(path), DeviceAdded, dev}
				emitter.Emit("discovery", devInfo)

				break
			}

		}
	})()

	return nil
}
