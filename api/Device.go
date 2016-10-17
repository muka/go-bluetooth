package api

import (
	"github.com/muka/bluez-client/bluez"
)

// NewDevice creates a new Device
func NewDevice(path string, props *bluez.Device1Properties) *Device {

	d := new(Device)

	d.Path = path
	d.Properties = props

	d.client = bluez.NewDevice1(path)
	d.client.Properties = props
	return d
}

//Device return an API to interact with a DBus device
type Device struct {
	Path       string
	Properties *bluez.Device1Properties
	client     *bluez.Device1
}

//GetClient return a DBus Device1 interface client
func (d *Device) GetClient() *bluez.Device1 {
	return d.client
}
