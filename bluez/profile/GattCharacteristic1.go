package profile

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewGattCharacteristic1 create a new GattCharacteristic1 client
func NewGattCharacteristic1(path string) (*GattCharacteristic1, error) {
	g := new(GattCharacteristic1)
	g.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: bluez.GattCharacteristic1Interface,
			Path:  path,
			Bus:   bluez.SystemBus,
		},
	)

	g.Path = path
	g.Properties = new(GattCharacteristic1Properties)

	_, err := g.GetProperties()
	if err != nil {
		return nil, err
	}

	return g, nil
}

// GattCharacteristic1 client
type GattCharacteristic1 struct {
	client     *bluez.Client
	Properties *GattCharacteristic1Properties
	channel    chan *dbus.Signal
	Path       string
}

// GattCharacteristic1Properties exposed properties for GattCharacteristic1
type GattCharacteristic1Properties struct {
	Value          []byte `dbus:"emit"`
	Notifying      bool
	NotifyAcquired bool
	WriteAcquired  bool
	Service        dbus.ObjectPath
	UUID           string
	Flags          []string
	Descriptors    []dbus.ObjectPath
}

//ToMap serialize properties
func (d *GattCharacteristic1Properties) ToMap() (map[string]interface{}, error) {
	if !d.Service.IsValid() {
		return nil, errors.New("GattCharacteristic1Properties: Service ObjectPath is not valid")
	}
	for i := 0; i < len(d.Descriptors); i++ {
		if d.Descriptors[i].IsValid() {
			return nil, errors.New("GattCharacteristic1Properties: Descriptors contains an ObjectPath that is not valid")
		}
	}
	return structs.Map(d), nil
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
func (d *GattCharacteristic1) Unregister(signal chan *dbus.Signal) error {
	if d.channel != nil {
		close(d.channel)
	}
	return d.client.Unregister(d.client.Config.Path, bluez.PropertiesInterface, signal)
}

//GetProperties load all available properties
func (d *GattCharacteristic1) GetProperties() (*GattCharacteristic1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//GetProperty load a single property
func (d *GattCharacteristic1) GetProperty(name string) (interface{}, error) {
	val, err := d.client.GetProperty(name)
	if err != nil {
		return nil, err
	}
	return val.Value(), nil
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

//AcquireWrite acquire file descriptor and MTU for writing [experimental]
func (d *GattCharacteristic1) AcquireWrite() (dbus.UnixFD, uint16, error) {
	var fd dbus.UnixFD
	var mtu uint16
	err := d.client.Call("AcquireWrite", 0).Store(&fd, &mtu)
	return fd, mtu, err
}

//AcquireNotify acquire file descriptor and MTU for notify [experimental]
func (d *GattCharacteristic1) AcquireNotify() (dbus.UnixFD, uint16, error) {
	var fd dbus.UnixFD
	var mtu uint16
	err := d.client.Call("AcquireNotify", 0).Store(&fd, &mtu)
	return fd, mtu, err
}
