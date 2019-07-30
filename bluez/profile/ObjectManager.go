package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewObjectManager create a new ObjectManager client
func NewObjectManager(name string, path string) (*ObjectManager, error) {
	om := new(ObjectManager)
	om.client = bluez.NewClient(
		&bluez.Config{
			Name:  name,
			Iface: "org.freedesktop.DBus.ObjectManager",
			Path:  path,
			Bus:   bluez.SystemBus,
		},
	)
	return om, nil
}

// ObjectManager manges the list of all available objects
type ObjectManager struct {
	client *bluez.Client
}

// Close the connection
func (o *ObjectManager) Close() {
	o.client.Disconnect()
}

// GetManagedObjects return a list of all available objects registered
func (o *ObjectManager) GetManagedObjects() (map[dbus.ObjectPath]map[string]map[string]dbus.Variant, error) {
	var objs map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err := o.client.Call("GetManagedObjects", 0).Store(&objs)
	return objs, err
}

//Register watch for signal events
func (o *ObjectManager) Register() (chan *dbus.Signal, error) {
	path := o.client.Config.Path
	iface := o.client.Config.Iface
	return o.client.Register(path, iface)
}

//Unregister watch for signal events
func (o *ObjectManager) Unregister(signal chan *dbus.Signal) error {
	path := o.client.Config.Path
	iface := o.client.Config.Iface
	return o.client.Unregister(path, iface, signal)
}

// SignalAdded notify of interfaces being added
func (o *ObjectManager) SignalAdded(path dbus.ObjectPath, props map[string]map[string]dbus.Variant) error {
	return o.client.Emit(path, bluez.InterfacesAdded, props)
}

// SignalRemoved notify of interfaces being removed
func (o *ObjectManager) SignalRemoved(path dbus.ObjectPath, ifaces []string) error {
	if ifaces == nil {
		ifaces = make([]string, 0)
	}
	return o.client.Emit(path, bluez.InterfacesRemoved, ifaces)
}

// GetManagedObject return an up to date view of a single object state.
// object is nil if the object path is not found
func (o *ObjectManager) GetManagedObject(objpath dbus.ObjectPath) (map[string]map[string]dbus.Variant, error) {
	objects, err := o.GetManagedObjects()
	if err != nil {
		return nil, err
	}
	if p, ok := objects[objpath]; ok {
		return p, nil
	}
	// return nil, errors.New("Object not found")
	return nil, nil
}

// AddObject an object to the managed list
func (o *ObjectManager) AddObject(path dbus.ObjectPath, object map[string]bluez.Properties) error {

	// 			 ifaces   	prop name    prop value
	obj := map[string]map[string]dbus.Variant{}
	for iface, props := range object {

		if _, ok := obj[iface]; !ok {
			obj[iface] = make(map[string]dbus.Variant)
		}

		propsVal, err := props.ToMap()
		if err != nil {
			return err
		}

		for key, value := range propsVal {
			obj[iface][key] = dbus.MakeVariant(value)
		}
	}

	return o.SignalAdded(path, obj)
}

// RemoveObject an object from the managed list
func (o *ObjectManager) RemoveObject(path dbus.ObjectPath) error {
	objects, err := o.GetManagedObjects()
	if err != nil {
		return err
	}

	if s, ok := objects[path]; ok {
		delete(objects, path)
		ifaces := make([]string, len(s))
		for i := range s {
			ifaces = append(ifaces, i)
		}
		return o.SignalRemoved(path, ifaces)
	}

	return nil
}
