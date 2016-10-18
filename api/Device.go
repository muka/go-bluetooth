package api

import (
	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/device-manager/util"
)

// NewDevice creates a new Device
func NewDevice(path string) *Device {
	d := new(Device)
	d.Path = path
	d.client = bluez.NewDevice1(path)
	return d
}

// ParseDevice parse a Device from a ObjectManager map
func ParseDevice(path dbus.ObjectPath, propsMap map[string]dbus.Variant) *Device {

	d := new(Device)
	d.Path = string(path)
	d.client = bluez.NewDevice1(d.Path)

	props := new(bluez.Device1Properties)
	util.MapToStruct(props, propsMap)
	d.client.Properties = props

	return d
}

//Device return an API to interact with a DBus device
type Device struct {
	Path string
	// Properties *bluez.Device1Properties
	client *bluez.Device1
}

//GetClient return a DBus Device1 interface client
func (d *Device) GetClient() *bluez.Device1 {
	return d.client
}

//GetProperties return the properties for the device
func (d *Device) GetProperties() *bluez.Device1Properties {
	if d.client.Properties == nil {
		d.client.GetProperties()
	}
	return d.client.Properties
}
