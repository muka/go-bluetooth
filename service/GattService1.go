package service

import (
	"strconv"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewGattService1 create a new instance of GattService1
func NewGattService1(config *GattService1Config, props *profile.GattService1Properties) *GattService1 {
	s := &GattService1{
		config:     config,
		properties: props,
	}
	return s
}

//GattService1Config GattService configuration
type GattService1Config struct {
	basePath   dbus.ObjectPath
	app        *Application
	ID         int
	objectPath dbus.ObjectPath
}

//GattService1 interface implementation
type GattService1 struct {
	config          *GattService1Config
	properties      *profile.GattService1Properties
	characteristics map[dbus.ObjectPath]*GattCharacteristic1
	charIndex       int
}

//GetApp return the parent app
func (s *GattService1) GetApp() *Application {
	return s.config.app
}

//Path return the object path
func (s *GattService1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Properties return the properties of the service
func (s *GattService1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	s.properties.Characteristics = s.GetCharacteristicPaths()
	p[bluez.GattService1Interface] = s.properties
	return p
}

//GetCharacteristics return the characteristics of the service
func (s *GattService1) GetCharacteristics() []*GattCharacteristic1 {
	chars := make([]*GattCharacteristic1, 0)
	for _, char := range s.characteristics {
		chars = append(chars, char)
	}
	return chars
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
	path := string(s.config.objectPath) + "/char" + strconv.Itoa(s.charIndex)
	config := &GattCharacteristic1Config{
		ID:         s.charIndex,
		objectPath: dbus.ObjectPath(path),
	}
	s.charIndex++
	char := NewGattCharacteristic1(config, props)
	return char
}

//AddCharacteristic add a characteristic
func (s *GattService1) AddCharacteristic(char *GattCharacteristic1) {
	s.characteristics[char.Path()] = char
	s.config.app.objectManager.AddObject(char.Path(), char.Properties())
}

//RemoveCharacteristic remove a characteristic
func (s *GattService1) RemoveCharacteristic(char *GattCharacteristic1) {
	if _, ok := s.characteristics[char.Path()]; ok {
		delete(s.characteristics, char.Path())
		s.config.app.objectManager.RemoveObject(char.Path())
	}
}
