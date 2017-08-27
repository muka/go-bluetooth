package profile

import (
	"github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewGattService1 create a new GattService1 client
func NewGattService1(path string, name string) *GattService1 {
	a := new(GattService1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  name,
			Iface: "org.bluez.GattService1",
			Path:  path,
			Bus:   bluez.SystemBus,
		},
	)
	a.Properties = new(GattService1Properties)
	a.GetProperties()
	return a
}

// GattService1 client
type GattService1 struct {
	client     *bluez.Client
	Properties *GattService1Properties
}

// GattService1Properties exposed properties for GattService1
type GattService1Properties struct {
	Primary         bool
	Device          dbus.ObjectPath
	Characteristics []dbus.ObjectPath `dbus:"emit"`
	UUID            string
}

//ToMap serialize a properties struct to a map
func (p *GattService1Properties) ToMap() (map[string]interface{}, error) {

	m := structs.Map(p)

	if !p.Device.IsValid() {
		delete(m, "Device")
	}

	chars := make([]dbus.ObjectPath, 0)
	for _, c := range p.Characteristics {
		if c.IsValid() {
			chars = append(chars, c)
		}
	}
	m["Characteristics"] = chars

	return m, nil
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
