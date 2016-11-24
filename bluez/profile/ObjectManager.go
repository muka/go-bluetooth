package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewObjectManager create a new Device1 client
func NewObjectManager() *ObjectManager {
	om := new(ObjectManager)
	om.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: "org.freedesktop.DBus.ObjectManager",
			Path:  "/",
			Bus:   bluez.SystemBus,
		},
	)

	return om
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
func (o *ObjectManager) Unregister() error {
	path := o.client.Config.Path
	iface := o.client.Config.Iface
	return o.client.Unregister(path, iface)
}
