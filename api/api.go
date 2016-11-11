package api

import (
	"errors"
	"strings"

	"github.com/godbus/dbus"
	"github.com/op/go-logging"
	"github.com/tj/go-debug"

	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/bluez/profile"
	"github.com/muka/bluez-client/emitter"
	"github.com/muka/bluez-client/util"
)

var logger = logging.MustGetLogger("api")

var manager *profile.ObjectManager

var dbg = debug.Debug("bluez:api")

var registrations map[string]*dbus.Signal

var watchChangesEnabled = false

//Exit performs a clean exit
func Exit() {
	GetManager().Unregister()
	GetManager().Close()
	dbg("Bye.")
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

	var managerDbg = debug.Debug("bluez:api:manger-watch")

	if watchChangesEnabled {
		return nil
	}

	managerDbg("Watching manager changes")

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

			managerDbg("waiting for updates")

			if channel == nil {
				managerDbg("quitting manager listener")
				break
			}

			v := <-channel

			managerDbg("update received %s from %s", v.Name, v.Sender)

			if v == nil {
				managerDbg("nil value, exit loop")
				return
			}

			switch v.Name {
			case bluez.InterfacesAdded:
				{
					managerDbg("Added %s %s", v.Name, v.Path)

					// dbg("\n+++Body %s\n", v.Body)
					path := v.Body[0].(dbus.ObjectPath)
					props := v.Body[1].(map[string]map[string]dbus.Variant)

					managerDbg("Body %v", props)

					//device added
					if props[bluez.Device1Interface] != nil {
						dev, err := ParseDevice(path, props[bluez.Device1Interface])
						if err != nil {
							logger.Fatalf("Failed to parse device: %v\n", err)
							return
						}

						managerDbg("Added device %s", path)
						devInfo := DiscoveredDeviceEvent{string(path), DeviceAdded, dev}
						emitter.Emit("discovery", devInfo)
					}

					// Adapter added
					if props[bluez.Adapter1Interface] != nil {
						strpath := string(path)
						parts := strings.Split(strpath, "/")
						name := parts[len(parts)-1:][0]

						managerDbg("Added adapter %s", name)
						adapterInfo := AdapterEvent{name, strpath, DeviceAdded}
						emitter.Emit("adapter", adapterInfo)
					}

					//GattService1
					if props[bluez.GattService1Interface] != nil {

						strpath := string(path)
						parts := strings.Split(strpath, "/")
						devicePath := strings.Join(parts[:len(parts)-1], "/")

						managerDbg("Added GattService1 %s", strpath)

						srvcProps := new(profile.GattService1Properties)
						util.MapToStruct(srvcProps, props[bluez.GattService1Interface])

						ev := GattServiceEvent{strpath, devicePath, srvcProps, StatusAdded}

						emitter.Emit("service", ev)
						emitter.Emit(devicePath+".service", ev)

					}
					//GattCharacteristic1
					if props[bluez.GattCharacteristic1Interface] != nil {

						strpath := string(path)
						parts := strings.Split(strpath, "/")
						devicePath := strings.Join(parts[:len(parts)-2], "/")

						managerDbg("Added GattCharacteristic1 %s", strpath)

						srvcProps := new(profile.GattCharacteristic1Properties)
						util.MapToStruct(srvcProps, props[bluez.GattCharacteristic1Interface])

						ev := GattCharacteristicEvent{strpath, devicePath, srvcProps, StatusAdded}

						emitter.Emit("char", ev)
						emitter.Emit(devicePath+".char", ev)
					}
					//GattDescriptor1
					if props[bluez.GattDescriptor1Interface] != nil {
						strpath := string(path)
						parts := strings.Split(strpath, "/")
						devicePath := strings.Join(parts[:len(parts)-3], "/")

						managerDbg("Added GattDescriptor1 %s", strpath)

						srvcProps := new(profile.GattDescriptor1Properties)
						util.MapToStruct(srvcProps, props[bluez.GattDescriptor1Interface])

						ev := GattDescriptorEvent{strpath, devicePath, srvcProps, StatusAdded}

						emitter.Emit("char", ev)
						emitter.Emit(devicePath+".char", ev)
					}

				}
			case bluez.InterfacesRemoved:
				{

					managerDbg("Removed %s %s", v.Name, v.Path)

					// dbg("\n+++Body %s\n", v.Body)
					path := v.Body[0].(dbus.ObjectPath)
					ifaces := v.Body[1].([]string)
					for _, iF := range ifaces {
						// device removed
						if iF == bluez.Device1Interface {
							// dbg("%s : %s", path, ifaces)
							managerDbg("Removed device %s", path)
							devInfo := DiscoveredDeviceEvent{string(path), DeviceRemoved, nil}
							emitter.Emit("discovery", devInfo)
						}
						//adapter removed
						if iF == bluez.Adapter1Interface {
							// dbg("%s : %s", path, ifaces)
							strpath := string(path)
							parts := strings.Split(strpath, "/")
							name := parts[len(parts)-1:][0]

							managerDbg("Removed adapter %s", name)
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
