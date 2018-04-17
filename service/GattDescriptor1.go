package service

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewGattDescriptor1 create a new GattDescriptor1 client
func NewGattDescriptor1(config *GattDescriptor1Config, props *profile.GattDescriptor1Properties) (*GattDescriptor1, error) {

	propInterface, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	g := &GattDescriptor1{
		config:              config,
		properties:          props,
		PropertiesInterface: propInterface,
	}

	err = propInterface.AddProperties(g.Interface(), props)
	if err != nil {
		return nil, err
	}

	return g, nil
}

//GattDescriptor1Config GattDescriptor1 configuration
type GattDescriptor1Config struct {
	objectPath     dbus.ObjectPath
	characteristic *GattCharacteristic1
	ID             int
	conn           *dbus.Conn
}

// GattDescriptor1 client
type GattDescriptor1 struct {
	config              *GattDescriptor1Config
	properties          *profile.GattDescriptor1Properties
	PropertiesInterface *Properties
}

//Path return the object path
func (s *GattDescriptor1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Interface return the Dbus interface
func (s *GattDescriptor1) Interface() string {
	return bluez.GattDescriptor1Interface
}

//Properties return the properties of the service
func (s *GattDescriptor1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	s.properties.Characteristic = s.config.characteristic.Path()
	p[s.Interface()] = s.properties
	return p
}

//ReadValue read a value
func (s *GattDescriptor1) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {
	b, err := s.config.characteristic.config.service.config.app.HandleDescriptorRead(
		s.config.characteristic.config.service.properties.UUID, s.config.characteristic.properties.UUID,
		s.properties.UUID)

	var dberr *dbus.Error = nil

	if err != nil {
		if err.code == -1 {
			// No registered callback, so we'll just use our stored value
			b = s.properties.Value
		} else {
			dberr = dbus.NewError(err.Error(), nil)
		}
	}

	return b, dberr
}

//WriteValue write a value
func (s *GattDescriptor1) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {
	err := s.config.characteristic.config.service.config.app.HandleDescriptorWrite(
		s.config.characteristic.config.service.properties.UUID, s.config.characteristic.properties.UUID,
		s.properties.UUID, value)

	if err != nil {
		if err.code == -1 {
			// No registered callback, so we'll just store this value
			s.UpdateValue(value)
			return nil
		} else {
			dberr := dbus.NewError(err.Error(), nil)
			return dberr
		}
	}

	return nil
}

func (s *GattDescriptor1) UpdateValue(value []byte) {
	s.properties.Value = value
	s.PropertiesInterface.Instance().Set(s.Interface(), "Value", dbus.MakeVariant(value))
}

//Expose the desc to dbus
func (s *GattDescriptor1) Expose() error {

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
