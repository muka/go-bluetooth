package service

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

// NewObjectManagerService create an object manager to manage the exposed services of the app
func NewObjectManagerService(conn *dbus.Conn) (*ObjectManagerService, error) {

	o := &ObjectManagerService{
		conn:    conn,
		objects: make(map[dbus.ObjectPath]map[string]bluez.Properties),
	}

	return o, nil
}

// ObjectManagerService interface implementation
type ObjectManagerService struct {
	conn    *dbus.Conn
	objects map[dbus.ObjectPath]map[string]bluez.Properties
}

// SignalAdded notify of interfaces being added
func (o *ObjectManagerService) SignalAdded(path dbus.ObjectPath) error {

	props, err := o.GetManagedObject(path)
	if err != nil {
		return err
	}

	return o.conn.Emit(path, bluez.InterfacesAdded, props)
}

// SignalRemoved notify of interfaces being removed
func (o *ObjectManagerService) SignalRemoved(path dbus.ObjectPath, ifaces []string) error {
	if ifaces == nil {
		ifaces = make([]string, 0)
	}
	return o.conn.Emit(path, bluez.InterfacesRemoved, ifaces)
}

// GetManagedObject return an up to date view of a single object state
func (o *ObjectManagerService) GetManagedObject(objpath dbus.ObjectPath) (map[string]map[string]dbus.Variant, error) {
	props, err := o.GetManagedObjects()
	if err != nil {
		return nil, err
	}
	if p, ok := props[objpath]; ok {
		return p, nil
	}
	return nil, errors.New("Object not found")
}

// GetManagedObjects return an up to date view of the object state
func (o *ObjectManagerService) GetManagedObjects() (map[dbus.ObjectPath]map[string]map[string]dbus.Variant, *dbus.Error) {

	props := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	for path, ifs := range o.objects {
		if _, ok := props[path]; !ok {
			props[path] = make(map[string]map[string]dbus.Variant)
		}
		for i, m := range ifs {
			if _, ok := props[path][i]; !ok {
				props[path][i] = make(map[string]dbus.Variant)
			}
			l, err := m.ToMap()
			if err != nil {
				log.Errorf("Failed to serialize properties: %s", err.Error())
				return nil, DbusErr
			}
			for k, v := range l {
				vrt := dbus.MakeVariant(v)
				props[path][i][k] = vrt
			}
		}
	}

	return props, nil
}

//AddObject add an object to the list
func (o *ObjectManagerService) AddObject(path dbus.ObjectPath, val map[string]bluez.Properties) error {
	o.objects[path] = val
	return o.SignalAdded(path)
}

//RemoveObject remove an object from the list
func (o *ObjectManagerService) RemoveObject(path dbus.ObjectPath) error {
	if s, ok := o.objects[path]; ok {
		delete(o.objects, path)
		ifaces := make([]string, len(s))
		for i := range s {
			ifaces = append(ifaces, i)
		}
		return o.SignalRemoved(path, ifaces)
	}
	return nil
}
