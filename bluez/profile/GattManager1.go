package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

// NewGattManager1 create a new GattManager1 client
func NewGattManager1(hostID string) *GattManager1 {
	a := new(GattManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.GattManager1",
			Path:  "/org/bluez/" + hostID,
			Bus:   bluez.SystemBus,
		},
	)
	a.Properties = new(GattManager1Properties)
	a.GetProperties()
	return a
}

// GattManager1 client
type GattManager1 struct {
	client     *bluez.Client
	Properties *GattManager1Properties
}

// GattManager1Properties exposed properties for GattManager1
type GattManager1Properties struct {
}

//GetProperties load all available properties
func (a *GattManager1) GetProperties() (*GattManager1Properties, error) {
	err := a.client.GetProperties(a.Properties)
	return a.Properties, err
}

// Close the connection
func (a *GattManager1) Close() {
	a.client.Disconnect()
}

//RegisterApplication add a new bluetooth Application
func (a *GattManager1) RegisterApplication(app dbus.ObjectPath, options map[string]interface{}) error {
	log.Debugf("Registering app %s", app)
	return a.client.Call("RegisterApplication", 0, app, options).Store()
}

//UnregisterApplication remove a bluetooth Application
func (a *GattManager1) UnregisterApplication(app dbus.ObjectPath) error {
	return a.client.Call("UnregisterApplication", 0, app).Store()
}
