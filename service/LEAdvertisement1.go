package service

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewLEAdvertisement1 create a new LEAdvertisement1 client
func NewLEAdvertisement1(config *LEAdvertisement1Config, props *profile.LEAdvertisement1Properties) (*LEAdvertisement1, error) {

	propInterface, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	s := &LEAdvertisement1{
		config:              config,
		properties:          props,
		PropertiesInterface: propInterface,
	}

	err = propInterface.AddProperties(s.Interface(), props)
	if err != nil {
		return nil, err
	}

	return s, nil
}

//LEAdvertisement1Config LEAdvertisement1 configuration
type LEAdvertisement1Config struct {
	objectPath dbus.ObjectPath
	conn       *dbus.Conn
}

// LEAdvertisement1 client
type LEAdvertisement1 struct {
	config              *LEAdvertisement1Config
	properties          *profile.LEAdvertisement1Properties
	PropertiesInterface *Properties
}

//Interface return the dbus interface name
func (s *LEAdvertisement1) Interface() string {
	return bluez.LEAdvertisement1Interface
}

//Path return the object path
func (s *LEAdvertisement1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Properties return the properties of the service
func (s *LEAdvertisement1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	p[s.Interface()] = s.properties
	return p
}

//Release This method gets called when the service daemon
// removes the Advertisement. A client can use it to do
// cleanup tasks. There is no need to call
// UnregisterAdvertisement because when this method gets
// called it has already been unregistered.
func (s *LEAdvertisement1) Release() {
	//callback here ?
}

//Expose the char to dbus
func (s *LEAdvertisement1) Expose() error {

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
			//LEAdvertisement1
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
