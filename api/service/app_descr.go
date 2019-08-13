package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

type Descr struct {
	app   *App
	path  dbus.ObjectPath
	props *gatt.GattDescriptor1Properties
}

func (s *Descr) Path() dbus.ObjectPath {
	return s.path
}

func (s *Descr) Interface() string {
	return gatt.GattDescriptor1Interface
}

func (s *Descr) Properties() bluez.Properties {
	return s.props
}

func (s *Descr) App() *App {
	return s.app
}

// Expose descr to dbus
func (s *Descr) Expose() error {
	return ExposeService(s)
}

// Remove descr from dbus
func (s *Descr) Remove() error {
	return RemoveService(s)
}
