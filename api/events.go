package api

import (
	"github.com/muka/bluez-client/bluez/profile"
	"github.com/muka/bluez-client/emitter"
)

//Event triggered
type Event emitter.Event

//Callback to be called on event
type Callback func(ev Event)

// DeviceStatus indicate the status of a device
type DeviceStatus int

const (
	//DeviceAdded indicates the device interface appeared
	DeviceAdded DeviceStatus = iota
	//DeviceRemoved indicates the device interface disappeared
	DeviceRemoved
)

//DiscoveredDeviceEvent contains detail regarding an ongoing discovery operation
type DiscoveredDeviceEvent struct {
	Path   string
	Status DeviceStatus
	Device *Device
}

// AdapterEvent reports the availability of a bluetooth adapter
type AdapterEvent struct {
	Name   string
	Path   string
	Status DeviceStatus
}

// PropertyChangedEvent an object to describe a changed property
type PropertyChangedEvent struct {
	Iface      string
	Field      string
	Value      interface{}
	Properties *profile.Device1Properties
}
