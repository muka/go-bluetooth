package service

import (
	"strconv"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewGattService1 create a new instance of GattService1
func NewGattService1(config *GattService1Config, props *profile.GattService1Properties) (*GattService1, error) {

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

	return s, nil
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
	properties          *profile.GattService1Properties
	characteristics     map[dbus.ObjectPath]*GattCharacteristic1
	charIndex           int
	PropertiesInterface *Properties
}

//Interface return the dbus interface name
func (s *GattService1) Interface() string {
	return bluez.GattService1Interface
}

//GetApp return the parent app
func (s *GattService1) GetApp() *Application {
	return s.config.app
}

//Path return the object path
func (s *GattService1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

func (s *GattService1) Advertised() bool {
	return s.config.advertised
}

//Properties return the properties of the service
func (s *GattService1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	s.properties.Characteristics = s.GetCharacteristicPaths()
	p[s.Interface()] = s.properties
	return p
}

//GetCharacteristics return the characteristics of the service
func (s *GattService1) GetCharacteristics() map[dbus.ObjectPath]*GattCharacteristic1 {
	return s.characteristics
}

//GetCharacteristicPaths return the characteristics object paths
func (s *GattService1) GetCharacteristicPaths() []dbus.ObjectPath {
	paths := make([]dbus.ObjectPath, 0)
	for path := range s.characteristics {
		paths = append(paths, path)
	}
	return paths
}

//CreateCharacteristic create a new characteristic
func (s *GattService1) CreateCharacteristic(props *profile.GattCharacteristic1Properties) (*GattCharacteristic1, error) {
	s.charIndex++
	path := string(s.config.objectPath) + "/char" + strconv.Itoa(s.charIndex)
	config := &GattCharacteristic1Config{
		ID:         s.charIndex,
		objectPath: dbus.ObjectPath(path),
		conn:       s.config.conn,
		service:    s,
	}

	props.Service = s.Path()

	char, err := NewGattCharacteristic1(config, props)
	return char, err
}

//AddCharacteristic add a characteristic
func (s *GattService1) AddCharacteristic(char *GattCharacteristic1) error {

	s.characteristics[char.Path()] = char

	err := char.Expose()
	if err != nil {
		return err
	}

	err = s.GetApp().exportTree()
	if err != nil {
		return err
	}

	om := s.config.app.GetObjectManager()
	return om.AddObject(char.Path(), char.Properties())
}

//RemoveCharacteristic remove a characteristic
func (s *GattService1) RemoveCharacteristic(char *GattCharacteristic1) error {
	if _, ok := s.characteristics[char.Path()]; ok {
		delete(s.characteristics, char.Path())
		om := s.config.app.GetObjectManager()
		return om.RemoveObject(char.Path())
	}
	return nil
}

//Expose the service to dbus
func (s *GattService1) Expose() error {

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
			//GattService1
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
