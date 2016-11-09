package api

import (
	"errors"
	"strings"

	"github.com/godbus/dbus"
	"github.com/juju/loggo"

	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/bluez/profile"
	"github.com/muka/bluez-client/emitter"
)

var logger = loggo.GetLogger("api")

var manager *profile.ObjectManager

var registrations map[string]*dbus.Signal

var watchChangesEnabled = false

//Exit performs a clean exit
func Exit() {
	GetManager().Unregister()
	GetManager().Close()
	logger.Infof("Bye.")
}

//GetManager return the object manager reference
func GetManager() *profile.ObjectManager {
	if manager == nil {
		manager = profile.NewObjectManager()
		WatchManagerChanges()
	}
	return manager
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

// UnWatchManagerChanges regitster for signals from the ObjectManager
func UnWatchManagerChanges() error {
	return GetManager().Unregister()
}

// WatchManagerChanges regitster for signals from the ObjectManager
func WatchManagerChanges() error {

	if watchChangesEnabled {
		return nil
	}

	m := GetManager()

	if m == nil {
		return nil
	}

	channel, err := m.Register()
	if err != nil {
		return err
	}

	// ensure is done once
	watchChangesEnabled = true

	go (func() {
		for {

			if channel == nil {
				logger.Println("Quitting manager listener")
				break
			}

			v := <-channel

			if v == nil {
				return
			}

			switch v.Name {
			case bluez.InterfacesAdded:
				{
					// logger.Printf("Added %s %s", v.Name, v.Path)
					// logger.Printf("\n+++Body %s\n", v.Body)
					path := v.Body[0].(dbus.ObjectPath)
					props := v.Body[1].(map[string]map[string]dbus.Variant)
					//device added
					if props[bluez.Device1Interface] != nil {
						dev, err := ParseDevice(path, props[bluez.Device1Interface])
						if err != nil {
							logger.Fatalf("Failed to parse device: %v\n", err)
							return
						}
						// logger.Printf("Added device %s", path)
						devInfo := DiscoveredDeviceEvent{string(path), DeviceAdded, dev}
						emitter.Emit("discovery", devInfo)
					}
					//adapter added
					if props[bluez.Adapter1Interface] != nil {
						strpath := string(path)
						parts := strings.Split(strpath, "/")
						name := parts[len(parts)-1:][0]
						// logger.Printf("Added adapter %s", name)
						adapterInfo := AdapterEvent{name, strpath, DeviceAdded}
						emitter.Emit("adapter", adapterInfo)
					}
				}
			case bluez.InterfacesRemoved:
				{
					// logger.Printf("Removed %s %s", v.Name, v.Path)
					// logger.Printf("\n+++Body %s\n", v.Body)
					path := v.Body[0].(dbus.ObjectPath)
					ifaces := v.Body[1].([]string)
					for _, iF := range ifaces {
						// device removed
						if iF == bluez.Device1Interface {
							// logger.Printf("%s : %s", path, ifaces)
							// logger.Printf("Removed device %s", path)
							devInfo := DiscoveredDeviceEvent{string(path), DeviceRemoved, nil}
							emitter.Emit("discovery", devInfo)
						}
						//adapter removed
						if iF == bluez.Adapter1Interface {
							// logger.Printf("%s : %s", path, ifaces)
							strpath := string(path)
							parts := strings.Split(strpath, "/")
							name := parts[len(parts)-1:][0]
							logger.Printf("Removed adapter %s", name)
							adapterInfo := AdapterEvent{name, strpath, DeviceRemoved}
							emitter.Emit("adapter", adapterInfo)
						}
					}
				}
			}
		}
	})()
	return nil
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
