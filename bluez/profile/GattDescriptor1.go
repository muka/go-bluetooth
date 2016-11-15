package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewGattDescriptor1 create a new GattDescriptor1 client
func NewGattDescriptor1(path string) *GattDescriptor1 {
	a := new(GattDescriptor1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: bluez.GattDescriptor1Interface,
			Path:  path,
			Bus:   bluez.SystemBus,
		},
	)
	a.Properties = new(GattDescriptor1Properties)
	a.GetProperties()
	return a
}

// GattDescriptor1 client
type GattDescriptor1 struct {
	client     *bluez.Client
	Properties *GattDescriptor1Properties
}

// GattDescriptor1Properties exposed properties for GattDescriptor1
type GattDescriptor1Properties struct {
	Value          []byte
	Characteristic dbus.ObjectPath
	UUID           string
}

// Close the connection
func (d *GattDescriptor1) Close() {
	d.client.Disconnect()
}

//Register for changes signalling
func (d *GattDescriptor1) Register() (chan *dbus.Signal, error) {
	return d.client.Register(d.client.Config.Path, bluez.PropertiesInterface)
}

//Unregister for changes signalling
func (d *GattDescriptor1) Unregister() error {
	return d.client.Unregister(d.client.Config.Path, bluez.PropertiesInterface)
}

//GetProperties load all available properties
func (d *GattDescriptor1) GetProperties() (*GattDescriptor1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//ReadValue read a value from a descriptor
func (d *GattDescriptor1) ReadValue(options map[string]dbus.Variant) ([]byte, error) {
	var b []byte
	err := d.client.Call("ReadValue", 0, options).Store(&b)
	return b, err
}

//WriteValue write a value to a characteristic
func (d *GattDescriptor1) WriteValue(b []byte, options map[string]dbus.Variant) error {
	err := d.client.Call("WriteValue", 0, b, options).Store()
	return err
}
