package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

// NewGattService1 create a new instance of GattService1
func NewGattService1(config *GattService1Config, props *gatt.GattService1Properties) (*GattService1, error) {

	propInterface, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	s := &GattService1{
		config:              config,
		properties:          props,
		PropertiesInterface: propInterface,
		characteristics:     make(map[dbus.ObjectPath]*GattCharacteristic1, 0),
	}

	err = propInterface.AddProperties(s.Interface(), props)
	if err != nil {
		return nil, err
	}

	props.Includes = append(props.Includes, config.objectPath)

	return s, nil
}

func NewGattService1Properties(uuid string) *gatt.GattService1Properties {
	return &gatt.GattService1Properties{
		IsService: true,
		Primary:   true,
		UUID:      uuid,
		Includes:  []dbus.ObjectPath{},
	}
}

//GattService1Config GattService configuration
type GattService1Config struct {
	app        *Application
	ID         int
	objectPath dbus.ObjectPath
	conn       *dbus.Conn
	advertised bool
}

//GattService1 interface implementation
type GattService1 struct {
	config              *GattService1Config
	properties          *gatt.GattService1Properties
	characteristics     map[dbus.ObjectPath]*GattCharacteristic1
	charIndex           int
	PropertiesInterface *Properties
}

//GetApp return the parent app
func (s *GattService1) GetApp() *Application {
	return s.config.app
}

//Conn return the dbus conn
func (s *GattService1) Conn() *dbus.Conn {
	return s.config.conn
}

//Interface return the dbus interface name
func (s *GattService1) Interface() string {
	return gatt.GattService1Interface
}

//Path return the object path
func (s *GattService1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Advertised indicate if the service has been advertised
func (s *GattService1) Advertised() bool {
	return s.config.advertised
}

//Properties return the properties of the service
func (s *GattService1) GetProperties() *gatt.GattService1Properties {
	return s.properties
}

//Properties return the properties of the service
func (s *GattService1) Properties() bluez.Properties {
	s.properties.Characteristics = s.GetCharacteristicPaths()
	return s.properties
}

//Expose the service to dbus
func (s *GattService1) Expose() error {
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
	// 		//GattService1
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
	// return nil
}
