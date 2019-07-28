package api

import (
	"strings"
	"sync"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/src/gen/profile/adapter"
	"github.com/muka/go-bluetooth/src/gen/profile/device"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
	"github.com/muka/go-bluetooth/util"
	log "github.com/sirupsen/logrus"
)

var manager *Manager

//GetManager return the object manager reference
func GetManager() (*Manager, error) {
	if manager == nil {
		m, err := NewManager()
		if err != nil {
			return nil, err
		}
		manager = m
	}
	return manager, nil
}

// NewManager creates a new manager instance
func NewManager() (*Manager, error) {
	m := new(Manager)
	m.objectManager = profile.NewObjectManager("org.bluez", "/")

	// m.objects = make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	m.objects = new(sync.Map)

	// watch for signaling from ObjectManager
	m.watchChanges()

	// Load initial object cache and emit events
	err := m.LoadObjects()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Manager track changes in the bluez dbus tree reflecting protocol updates
type Manager struct {
	objectManager       *profile.ObjectManager
	watchChangesEnabled bool
	objects             *sync.Map
	channel             chan *dbus.Signal
}

// unwatchChanges register for signals from the ObjectManager
func (m *Manager) unwatchChanges() error {
	if m.channel != nil {
		close(m.channel)
	}
	m.watchChangesEnabled = false
	return m.objectManager.Unregister(m.channel)
}

// watchChanges regitster for signals from the ObjectManager
func (m *Manager) watchChanges() error {

	if m.watchChangesEnabled {
		return nil
	}

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
		for v := range channel {

			if v == nil {
				m.watchChangesEnabled = false
				return
			}

			// log.Debugf("ObjectManager event: %++v", v)

			switch v.Name {
			case bluez.InterfacesAdded:
				{
					path := v.Body[0].(dbus.ObjectPath)
					props := v.Body[1].(map[string]map[string]dbus.Variant)

					// keep cache up to date
					m.objects.Store(path, props)

					emitChanges(path, props)
				}
			case bluez.InterfacesRemoved:
				{

					path := v.Body[0].(dbus.ObjectPath)
					ifaces := v.Body[1].([]string)

					// keep cache up to date
					m.objects.Delete(path)

					for _, iF := range ifaces {
						// device removed
						if iF == device.Device1Interface {

							devInfo := DiscoveredDeviceEvent{string(path), DeviceRemoved, nil}
							emitter.Emit("discovery", devInfo)
						}
						//adapter removed
						if iF == adapter.Adapter1Interface {

							strpath := string(path)
							parts := strings.Split(strpath, "/")
							name := parts[len(parts)-1:][0]

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
	if props[device.Device1Interface] != nil {
		dev, err := ParseDevice(path, props[device.Device1Interface])
		if err != nil {
			log.Fatalf("Failed to parse device: %v\n", err)
			return
		}

		devInfo := DiscoveredDeviceEvent{string(path), DeviceAdded, dev}
		emitter.Emit("discovery", devInfo)
	}

	//Adapter1
	if props[adapter.Adapter1Interface] != nil {
		strpath := string(path)
		parts := strings.Split(strpath, "/")
		name := parts[len(parts)-1:][0]

		adapterInfo := AdapterEvent{name, strpath, DeviceAdded}
		emitter.Emit("adapter", adapterInfo)
	}

	//GattService1
	if props[gatt.GattService1Interface] != nil {

		strpath := string(path)
		parts := strings.Split(strpath, "/")
		devicePath := strings.Join(parts[:len(parts)-1], "/")

		srvcProps := new(gatt.GattService1Properties)
		util.MapToStruct(srvcProps, props[gatt.GattService1Interface])

		ev := GattServiceEvent{strpath, devicePath, srvcProps, StatusAdded}

		emitter.Emit("service", ev)
		emitter.Emit(devicePath+".service", ev)

	}
	//GattCharacteristic1
	if props[gatt.GattCharacteristic1Interface] != nil {

		strpath := string(path)
		parts := strings.Split(strpath, "/")
		devicePath := strings.Join(parts[:len(parts)-2], "/")

		srvcProps := new(gatt.GattCharacteristic1Properties)
		util.MapToStruct(srvcProps, props[gatt.GattCharacteristic1Interface])

		ev := GattCharacteristicEvent{strpath, devicePath, srvcProps, StatusAdded}

		emitter.Emit("char", ev)
		emitter.Emit(devicePath+".char", ev)
	}
	//GattDescriptor1
	if props[gatt.GattDescriptor1Interface] != nil {
		strpath := string(path)
		parts := strings.Split(strpath, "/")
		devicePath := strings.Join(parts[:len(parts)-3], "/")

		srvcProps := new(profile.GattDescriptor1Properties)
		util.MapToStruct(srvcProps, props[gatt.GattDescriptor1Interface])

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
	for path, object := range objs {
		m.objects.Store(path, object)
	}
	return nil
}

//GetObjects return the cached list of objects from the ObjectManager
func (m *Manager) GetObjects() *sync.Map {
	return m.objects
}

//RefreshState emit local manager objects and interfaces
func (m *Manager) RefreshState() error {

	err := m.LoadObjects()
	if err != nil {
		return err
	}

	objs := m.GetObjects()
	objs.Range(func(path, ifaces interface{}) bool {
		emitChanges(path.(dbus.ObjectPath), ifaces.(map[string]map[string]dbus.Variant))
		return true
	})

	return nil
}

//Close Close the Manager and free underlying resources
func (m *Manager) Close() {
	m.objectManager.Unregister(m.channel)
	m.objectManager.Close()
	m.objectManager = nil
}
