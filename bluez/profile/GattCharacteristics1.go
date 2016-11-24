package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewGattCharacteristic1 create a new GattCharacteristic1 client
func NewGattCharacteristic1(path string) *GattCharacteristic1 {
	g := new(GattCharacteristic1)
	g.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: bluez.GattCharacteristic1Interface,
			Path:  path,
			Bus:   bluez.SystemBus,
		},
	)

	g.Properties = new(GattCharacteristic1Properties)

	_, err := g.GetProperties()
	if err != nil {
		panic(err)
	}

	return g
}

// GattCharacteristic1 client
type GattCharacteristic1 struct {
	client     *bluez.Client
	Properties *GattCharacteristic1Properties
	channel    chan *dbus.Signal
}

// GattCharacteristic1Properties exposed properties for GattCharacteristic1
type GattCharacteristic1Properties struct {
	Value     []byte
	Flags     []string
	Notifying bool
	Service   dbus.ObjectPath
	UUID      string
}

// Close the connection
func (d *GattCharacteristic1) Close() {
	d.client.Disconnect()
}

//Register for changes signalling
func (d *GattCharacteristic1) Register() (chan *dbus.Signal, error) {
	if d.channel == nil {
		channel, err := d.client.Register(d.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		d.channel = channel
	}
	return d.channel, nil
}

//Unregister for changes signalling
func (d *GattCharacteristic1) Unregister() error {
	if d.channel != nil {
		close(d.channel)
	}
	return d.client.Unregister(d.client.Config.Path, bluez.PropertiesInterface)
}

//GetProperties load all available properties
func (d *GattCharacteristic1) GetProperties() (*GattCharacteristic1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//ReadValue read a value from a characteristic
func (d *GattCharacteristic1) ReadValue(options map[string]dbus.Variant) ([]byte, error) {
	var b []byte
	err := d.client.Call("ReadValue", 0, options).Store(&b)
	return b, err
}

//WriteValue write a value to a characteristic
func (d *GattCharacteristic1) WriteValue(b []byte, options map[string]dbus.Variant) error {
	err := d.client.Call("WriteValue", 0, b, options).Store()
	return err
}

//StartNotify start notifications
func (d *GattCharacteristic1) StartNotify() error {
	return d.client.Call("StartNotify", 0).Store()
}

//StopNotify stop notifications
func (d *GattCharacteristic1) StopNotify() error {
	return d.client.Call("StopNotify", 0).Store()
}
