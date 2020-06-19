package bluez

import (
	"github.com/godbus/dbus/v5"
)

var objectManager *ObjectManager

// GetObjectManager return a client instance of the Bluez object manager
func GetObjectManager() (*ObjectManager, error) {
	if objectManager != nil {
		return objectManager, nil
	}

	om, err := NewObjectManager(OrgBluezInterface, "/")
	if err != nil {
		return nil, err
	}

	objectManager = om
	return om, nil
}

// NewObjectManager create a new ObjectManager client
func NewObjectManager(name string, path string) (*ObjectManager, error) {
	om := new(ObjectManager)
	om.client = NewClient(
		&Config{
			Name:  name,
			Iface: "org.freedesktop.DBus.ObjectManager",
			Path:  dbus.ObjectPath(path),
			Bus:   SystemBus,
		},
	)
	return om, nil
}

// ObjectManager manges the list of all available objects
type ObjectManager struct {
	client *Client
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
	return o.client.Register(dbus.ObjectPath(path), iface)
}

//Unregister watch for signal events
func (o *ObjectManager) Unregister(signal chan *dbus.Signal) error {
	path := o.client.Config.Path
	iface := o.client.Config.Iface
	return o.client.Unregister(dbus.ObjectPath(path), iface, signal)
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
