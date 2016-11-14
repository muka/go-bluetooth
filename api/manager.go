package api

import (
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/bluez/profile"
	"github.com/muka/bluez-client/emitter"
	"github.com/muka/bluez-client/util"
	"github.com/tj/go-debug"
)

var managerDbg = debug.Debug("bluez:api:manger-watch")

var manager *profile.ObjectManager
var watchChangesEnabled = false

//GetManager return the object manager reference
func GetManager() *profile.ObjectManager {
	if manager == nil {
		manager = profile.NewObjectManager()
		WatchManagerChanges()
	}
	return manager
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
		managerDbg("waiting for updates")
		for v := range channel {

			// if channel == nil {
			// 	managerDbg("quitting manager listener")
			// 	watchChangesEnabled = false
			// 	break
			// }
			// v := <-channel

			if v == nil {
				managerDbg("nil value, abort")
				watchChangesEnabled = false
				return
			}

			managerDbg("update received %s from %s", v.Name, v.Sender)

			switch v.Name {
			case bluez.InterfacesAdded:
				{
					managerDbg("Added %s %s", v.Name, v.Path)

					// dbg("\n+++Body %s\n", v.Body)
					path := v.Body[0].(dbus.ObjectPath)
					props := v.Body[1].(map[string]map[string]dbus.Variant)

					managerDbg("Body %v", props)
					loadManagedObject(path, props)
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

func loadManagedObject(path dbus.ObjectPath, props map[string]map[string]dbus.Variant) {

	//Device1
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

	//Adapter1
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

		emitter.Emit("desc", ev)
		emitter.Emit(devicePath+".desc", ev)
	}

}

//RefreshManagerState emit local manager objects and interfaces
func RefreshManagerState() error {

	objs, err := GetManager().GetManagedObjects()
	if err != nil {
		return err
	}

	for path, ifaces := range objs {
		managerDbg("Managed %s: %v", path, ifaces)
		loadManagedObject(path, ifaces)
	}

	return nil
}
