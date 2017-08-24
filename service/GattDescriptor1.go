package service

import (
	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewGattDescriptor1 create a new GattDescriptor1 client
func NewGattDescriptor1(config *GattDescriptor1Config, props *profile.GattDescriptor1Properties) *GattDescriptor1 {
	g := &GattDescriptor1{
		config:     config,
		properties: props,
	}
	return g
}

//GattDescriptor1Config GattDescriptor1 configuration
type GattDescriptor1Config struct {
	objectPath     dbus.ObjectPath
	characteristic *GattCharacteristic1
	ID             int
}

// GattDescriptor1 client
type GattDescriptor1 struct {
	config     *GattDescriptor1Config
	properties *profile.GattDescriptor1Properties
}

//Path return the object path
func (s *GattDescriptor1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Properties return the properties of the service
func (s *GattDescriptor1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	p[bluez.GattDescriptor1Interface] = s.properties
	return p
}

//ReadValue read a value
func (s *GattDescriptor1) ReadValue(options map[string]interface{}) []byte {
	log.Debug("Descriptor.ReadValue")
	b := make([]byte, 0)
	return b
}

//WriteValue write a value
func (s *GattDescriptor1) WriteValue(value []byte, options map[string]interface{}) {
	log.Debug("Descriptor.ReadValue")
}
