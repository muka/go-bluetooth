package service

import (
	"strconv"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/prometheus/common/log"
)

// NewGattService1 create a new instance of GattService1
func NewGattService1(config *GattService1Config, props *profile.GattService1Properties) (*GattService1, error) {

	propInterface, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	s := &GattService1{
		config:              config,
		props:               props,
		PropertiesInterface: propInterface,
	}

	err = propInterface.AddProperties(bluez.GattService1Interface, props)
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
}

//GattService1 interface implementation
type GattService1 struct {
	config              *GattService1Config
	props               *profile.GattService1Properties
	characteristics     map[dbus.ObjectPath]*GattCharacteristic1
	charIndex           int
	PropertiesInterface *Properties
}

//GetApp return the parent app
func (s *GattService1) GetApp() *Application {
	return s.config.app
}

//Path return the object path
func (s *GattService1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Iface return the Dbus interface
func (s *GattService1) Iface() string {
	return bluez.GattService1Interface
}

//Properties return the properties of the service
func (s *GattService1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	s.props.Characteristics = s.GetCharacteristicPaths()
	p[bluez.GattService1Interface] = s.props
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
func (s *GattService1) CreateCharacteristic(props *profile.GattCharacteristic1Properties) *GattCharacteristic1 {
	s.charIndex++
	path := string(s.config.objectPath) + "/char" + strconv.Itoa(s.charIndex)
	config := &GattCharacteristic1Config{
		ID:         s.charIndex,
		objectPath: dbus.ObjectPath(path),
	}
	char := NewGattCharacteristic1(config, props)
	return char
}

//AddCharacteristic add a characteristic
func (s *GattService1) AddCharacteristic(char *GattCharacteristic1) error {
	s.characteristics[char.Path()] = char
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

	log.Debugf("GATT Service path %s", s.Path())
	conn := s.config.conn
	err := conn.Export(s, s.Path(), bluez.GattService1Interface)
	if err != nil {
		return err
	}

	serviceNode := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,

			{
				Name:    bluez.GattService1Interface,
				Methods: introspect.Methods(s),
			},
		},
	}

	err = conn.Export(
		introspect.NewIntrospectable(serviceNode),
		s.Path(),
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	s.PropertiesInterface.Expose(s.Path())

	log.Debugf("Exposed GATT service %s", s.Path())

	return nil
}
