package profile

import (
	"log"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/util"
)

// NewGattService1 create a new GattService1 client
func NewGattService1(path string) *GattService1 {
	a := new(GattService1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.GattService1",
			Path:  path,
			Bus:   bluez.SystemBus,
		},
	)
	a.logger = util.NewLogger(path)
	a.Properties = new(GattService1Properties)
	return a
}

// GattService1 client
type GattService1 struct {
	client     *bluez.Client
	logger     *log.Logger
	Properties *GattService1Properties
}

// GattService1Properties exposed properties for GattService1
type GattService1Properties struct {
	Characteristics []dbus.ObjectPath
	Primary         bool
	Device          dbus.ObjectPath
	UUID            string
}

// Close the connection
func (d *GattService1) Close() {
	d.client.Disconnect()
}

//Register for changes signalling
func (d *GattService1) Register() (chan *dbus.Signal, error) {
	return d.client.Register(d.client.Config.Path, bluez.PropertiesInterface)
}

//Unregister for changes signalling
func (d *GattService1) Unregister() error {
	return d.client.Unregister(d.client.Config.Path, bluez.PropertiesInterface)
}

//GetProperties load all available properties
func (d *GattService1) GetProperties() (*GattService1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}
