package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

// NewGattDescriptor1 create a new GattDescriptor1 client
func NewGattDescriptor1(config *GattDescriptor1Config, props *gatt.GattDescriptor1Properties) (*GattDescriptor1, error) {

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
	properties          *gatt.GattDescriptor1Properties
	PropertiesInterface *Properties
}

//Conn return a conn instance
func (s *GattDescriptor1) Conn() *dbus.Conn {
	return s.config.conn
}

//Path return the object path
func (s *GattDescriptor1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Interface return the Dbus interface
func (s *GattDescriptor1) Interface() string {
	return gatt.GattDescriptor1Interface
}

//Properties return the properties of the service
func (s *GattDescriptor1) Properties() bluez.Properties {
	s.properties.Characteristic = s.config.characteristic.Path()
	return s.properties
}

//Expose the desc to dbus
func (s *GattDescriptor1) Expose() error {
	return ExposeService(s)
	//
	// conn := s.config.conn
	//
	// err := conn.Export(s, s.Path(), s.Interface())
	// if err != nil {
	// 	return err
	// }
	//
	// for iface, props := range s.Properties() {
	// 	s.PropertiesInterface.AddProperties(iface, props)
	// }
	//
	// s.PropertiesInterface.Expose(s.Path())
	//
	// node := &introspect.Node{
	// 	Interfaces: []introspect.Interface{
	// 		//Introspect
	// 		introspect.IntrospectData,
	// 		//Properties
	// 		prop.IntrospectData,
	// 		//GattCharacteristic1
	// 		{
	// 			Name:       s.Interface(),
	// 			Methods:    introspect.Methods(s),
	// 			Properties: s.PropertiesInterface.Introspection(s.Interface()),
	// 		},
	// 	},
	// }
	//
	// err = conn.Export(
	// 	introspect.NewIntrospectable(node),
	// 	s.Path(),
	// 	"org.freedesktop.DBus.Introspectable")
	// if err != nil {
	// 	return err
	// }
	//
	// return nil
}
