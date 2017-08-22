package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewGattService1 create a new instance of GattService1
func NewGattService1(conn *dbus.Conn, basePath string, props *profile.GattService1Properties) *GattService1 {
	s := &GattService1{
		config: &GattService1Config{
			conn:     conn,
			basePath: basePath,
		},
		properties: props,
	}

	return s
}

//GattService1Config GattService configuration
type GattService1Config struct {
	basePath string
	conn     *dbus.Conn
}

//GattService1 interface implementation
type GattService1 struct {
	config     *GattService1Config
	objectPath dbus.ObjectPath
	properties bluez.Properties
}

//Path return the object path
func (s *GattService1) Path() dbus.ObjectPath {
	return s.objectPath
}

//Properties return the properties of the service
func (s *GattService1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	p[bluez.GattService1Interface] = s.properties
	return p
}
