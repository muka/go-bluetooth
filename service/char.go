package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
	log "github.com/sirupsen/logrus"
)

// NewGattCharacteristic1 create a new GattCharacteristic1 client
func NewGattCharacteristic1(config *GattCharacteristic1Config, props *gatt.GattCharacteristic1Properties) (*GattCharacteristic1, error) {

	propInterface, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	s := &GattCharacteristic1{
		config:              config,
		properties:          props,
		PropertiesInterface: propInterface,
		descriptors:         make(map[dbus.ObjectPath]*GattDescriptor1),
	}

	err = propInterface.AddProperties(s.Interface(), props)
	if err != nil {
		return nil, err
	}

	return s, nil
}

//GattCharacteristic1Config GattCharacteristic1 configuration
type GattCharacteristic1Config struct {
	objectPath dbus.ObjectPath
	service    *GattService1
	ID         int
	conn       *dbus.Conn
}

// GattCharacteristic1 client
type GattCharacteristic1 struct {
	config              *GattCharacteristic1Config
	properties          *gatt.GattCharacteristic1Properties
	PropertiesInterface *Properties
	descriptors         map[dbus.ObjectPath]*GattDescriptor1
	descIndex           int
	notifying           bool
}

//Interface return the dbus interface name
func (s *GattCharacteristic1) Conn() *dbus.Conn {
	return s.config.conn
}

//Interface return the dbus interface name
func (s *GattCharacteristic1) Interface() string {
	return gatt.GattCharacteristic1Interface
}

//Path return the object path
func (s *GattCharacteristic1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Properties return the properties of the service
func (s *GattCharacteristic1) Properties() bluez.Properties {
	s.properties.Descriptors = s.GetDescriptorPaths()
	return s.properties
}

//StartNotify start notification
func (s *GattCharacteristic1) StartNotify() *dbus.Error {
	log.Debug("Characteristic.StartNotify")
	s.notifying = true
	return nil
}

//StopNotify stop notification
func (s *GattCharacteristic1) StopNotify() *dbus.Error {
	log.Debug("Characteristic.StopNotify")
	s.notifying = false
	return nil
}

//Expose the char to dbus
func (s *GattCharacteristic1) Expose() error {
	return ExposeService(s)
}
