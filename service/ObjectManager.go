package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewObjectManager create a new instance
func NewObjectManager(conn *dbus.Conn) (*ObjectManager, error) {

	o := &ObjectManager{
		conn:    conn,
		objects: make(map[dbus.ObjectPath]map[string]bluez.Properties),
	}

	return o, nil
}

// ObjectManager interface implementation
type ObjectManager struct {
	conn    *dbus.Conn
	objects map[dbus.ObjectPath]map[string]bluez.Properties
}

// GetManagedObjects return an up to date view of the object state
func (o *ObjectManager) GetManagedObjects() (map[dbus.ObjectPath]map[string]map[string]dbus.Variant, error) {

	props := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	for path, ifs := range o.objects {
		if _, ok := props[path]; !ok {
			props[path] = make(map[string]map[string]dbus.Variant)
		}
		for i, m := range ifs {
			if _, ok := props[path][i]; !ok {
				props[path][i] = make(map[string]dbus.Variant)
			}
			l := m.ToMap()
			for k, v := range l {
				props[path][i][k] = dbus.MakeVariant(v)
			}
		}
	}

	return props, nil
}

//AddObject add an object to the list
func (o *ObjectManager) AddObject(path dbus.ObjectPath, val map[string]bluez.Properties) error {
	o.objects[path] = val

	// signal to the bus
	o.conn.BusObject().Call(bluez.InterfacesAdded, 0,
		"type='signal',path='/org/freedesktop/DBus',interface='org.freedesktop.DBus',sender='org.freedesktop.DBus'")

	//TODO: handle dbus.Call error

	return nil
}

//RemoveObject remove an object from the list
func (o *ObjectManager) RemoveObject(path dbus.ObjectPath) error {
	if _, ok := o.objects[path]; ok {
		delete(o.objects, path)
		// signal to the bus
		o.conn.BusObject().Call(bluez.InterfacesRemoved, 0,
			"type='signal',path='/org/freedesktop/DBus',interface='org.freedesktop.DBus',sender='org.freedesktop.DBus'")
		//TODO: handle dbus.Call error
	}
	return nil
}
