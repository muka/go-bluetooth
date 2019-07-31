package service

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
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
func (s *GattCharacteristic1) Interface() string {
	return gatt.GattCharacteristic1Interface
}

//Path return the object path
func (s *GattCharacteristic1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Properties return the properties of the service
func (s *GattCharacteristic1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	s.properties.Descriptors = s.GetDescriptorPaths()
	p[s.Interface()] = s.properties
	return p
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

	conn := s.config.conn

	err := conn.Export(s, s.Path(), s.Interface())
	if err != nil {
		return err
	}

	for iface, props := range s.Properties() {
		s.PropertiesInterface.AddProperties(iface, props)
	}

	s.PropertiesInterface.Expose(s.Path())

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//Properties
			prop.IntrospectData,
			//GattCharacteristic1
			{
				Name:       s.Interface(),
				Methods:    introspect.Methods(s),
				Properties: s.PropertiesInterface.Introspection(s.Interface()),
			},
		},
	}

	err = conn.Export(
		introspect.NewIntrospectable(node),
		s.Path(),
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	return nil
}
