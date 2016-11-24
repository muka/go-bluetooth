package api

import (
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/util"
	"github.com/tj/go-debug"
)

var dbgManager = debug.Debug("bluez:api:Manager")
var manager *Manager

//GetManager return the object manager reference
func GetManager() *Manager {
	if manager == nil {
		manager = NewManager()
	}
	return manager
}

// NewManager creates a new manager instance
func NewManager() *Manager {
	m := new(Manager)
	m.objectManager = profile.NewObjectManager()
	m.objects = make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)

	// watch for signaling from ObjectManager
	m.watchChanges()

	// Load initial object cache and emit events
	err := m.LoadObjects()
	if err != nil {
		panic(err)
	}

	dbgManager("Manager initialized")
	return m
}

// Manager track changes in the bluez dbus tree reflecting protocol updates
type Manager struct {
	objectManager       *profile.ObjectManager
	watchChangesEnabled bool
	objects             map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	channel             chan *dbus.Signal
}

// unwatchChanges register for signals from the ObjectManager
func (m *Manager) unwatchChanges() error {
	if m.channel != nil {
		close(m.channel)
	}
	m.watchChangesEnabled = false
	return m.objectManager.Unregister()
}

// watchChanges regitster for signals from the ObjectManager
func (m *Manager) watchChanges() error {

	if m.watchChangesEnabled {
		return nil
	}

	dbgManager("Watching manager changes")

	if m == nil {
		return nil
	}

	channel, err := m.objectManager.Register()
	if err != nil {
		return err
	}
	m.channel = channel

	// ensure is done once
	m.watchChangesEnabled = true

	go (func() {
		dbgManager("waiting updates")
		for v := range channel {

			if v == nil {
				dbgManager("nil value, abort")
				m.watchChangesEnabled = false
				return
			}

			dbgManager("update received %s from %s", v.Name, v.Sender)

			switch v.Name {
			case bluez.InterfacesAdded:
				{
					dbgManager("Added %s %s", v.Name, v.Path)

					path := v.Body[0].(dbus.ObjectPath)
					props := v.Body[1].(map[string]map[string]dbus.Variant)

					// keep cache up to date
					m.objects[path] = props

					dbgManager("Body %v", props)
					emitChanges(path, props)
				}
			case bluez.InterfacesRemoved:
				{

					dbgManager("Removed %s %s", v.Name, v.Path)

					// dbg("\n+++Body %s\n", v.Body)
					path := v.Body[0].(dbus.ObjectPath)
					ifaces := v.Body[1].([]string)

					// keep cache up to date
					if _, ok := m.objects[path]; ok {
						delete(m.objects, path)
					}

					for _, iF := range ifaces {
						// device removed
						if iF == bluez.Device1Interface {
							// dbg("%s : %s", path, ifaces)
							dbgManager("Removed device %s", path)
							devInfo := DiscoveredDeviceEvent{string(path), DeviceRemoved, nil}
							emitter.Emit("discovery", devInfo)
						}
						//adapter removed
						if iF == bluez.Adapter1Interface {
							// dbg("%s : %s", path, ifaces)
							strpath := string(path)
							parts := strings.Split(strpath, "/")
							name := parts[len(parts)-1:][0]

							dbgManager("Removed adapter %s", name)
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

func emitChanges(path dbus.ObjectPath, props map[string]map[string]dbus.Variant) {

	//Device1
	if props[bluez.Device1Interface] != nil {
		dev, err := ParseDevice(path, props[bluez.Device1Interface])
		if err != nil {
			logger.Fatalf("Failed to parse device: %v\n", err)
			return
		}
		dbgManager("Added device %s", path)
		devInfo := DiscoveredDeviceEvent{string(path), DeviceAdded, dev}
		emitter.Emit("discovery", devInfo)
	}

	//Adapter1
	if props[bluez.Adapter1Interface] != nil {
		strpath := string(path)
		parts := strings.Split(strpath, "/")
		name := parts[len(parts)-1:][0]

		dbgManager("Added adapter %s", name)
		adapterInfo := AdapterEvent{name, strpath, DeviceAdded}
		emitter.Emit("adapter", adapterInfo)
	}

	//GattService1
	if props[bluez.GattService1Interface] != nil {

		strpath := string(path)
		parts := strings.Split(strpath, "/")
		devicePath := strings.Join(parts[:len(parts)-1], "/")

		dbgManager("Added GattService1 %s", strpath)

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

		dbgManager("Added GattCharacteristic1 %s", strpath)

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

		dbgManager("Added GattDescriptor1 %s", strpath)

		srvcProps := new(profile.GattDescriptor1Properties)
		util.MapToStruct(srvcProps, props[bluez.GattDescriptor1Interface])

		ev := GattDescriptorEvent{strpath, devicePath, srvcProps, StatusAdded}

		emitter.Emit("desc", ev)
		emitter.Emit(devicePath+".desc", ev)
	}

}

//LoadObjects force reloading of cache objects list
func (m *Manager) LoadObjects() error {
	objs, err := m.objectManager.GetManagedObjects()
	if err != nil {
		return err
	}
	m.objects = objs
	dbg("Loaded %d objects", len(m.objects))
	return nil
}

//GetObjects return the cached list of objects from the ObjectManager
func (m *Manager) GetObjects() *map[dbus.ObjectPath]map[string]map[string]dbus.Variant {
	return &m.objects
}

//RefreshState emit local manager objects and interfaces
func (m *Manager) RefreshState() error {

	err := m.LoadObjects()
	if err != nil {
		return err
	}

	dbgManager("Refreshing object state")
	objs := m.GetObjects()
	for path, ifaces := range *objs {
		emitChanges(path, ifaces)
	}

	return nil
}

//Close Close the Manager and free underlying resources
func (m *Manager) Close() {
	m.objectManager.Unregister()
	m.objectManager.Close()
	m.objectManager = nil
}
