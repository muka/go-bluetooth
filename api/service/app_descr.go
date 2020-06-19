package service

import (
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

type DescrReadCallback func(c *Descr, options map[string]interface{}) ([]byte, error)
type DescrWriteCallback func(c *Descr, value []byte) ([]byte, error)

type Descr struct {
	UUID string
	app  *App
	char *Char

	path dbus.ObjectPath

	Properties *gatt.GattDescriptor1Properties
	iprops     *api.DBusProperties

	readCallback  DescrReadCallback
	writeCallback DescrWriteCallback
}

func (s *Descr) DBusProperties() *api.DBusProperties {
	return s.iprops
}

func (s *Descr) DBusObjectManager() *api.DBusObjectManager {
	return s.App().DBusObjectManager()
}

func (s *Descr) DBusConn() *dbus.Conn {
	return s.App().DBusConn()
}

func (s *Descr) Path() dbus.ObjectPath {
	return s.path
}

func (s *Descr) Char() *Char {
	return s.char
}

func (s *Descr) Interface() string {
	return gatt.GattDescriptor1Interface
}

func (s *Descr) GetProperties() bluez.Properties {
	s.Properties.Characteristic = s.char.Path()
	return s.Properties
}

func (s *Descr) App() *App {
	return s.app
}

// Expose descr to dbus
func (s *Descr) Expose() error {
	return api.ExposeDBusService(s)
}

// Remove descr from dbus
func (s *Descr) Remove() error {
	return api.RemoveDBusService(s)
}
