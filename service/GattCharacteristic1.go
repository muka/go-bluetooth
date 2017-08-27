package service

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

// NewGattCharacteristic1 create a new GattCharacteristic1 client
func NewGattCharacteristic1(config *GattCharacteristic1Config, props *profile.GattCharacteristic1Properties) *GattCharacteristic1 {
	g := &GattCharacteristic1{
		config:     config,
		properties: props,
	}
	return g
}

//GattCharacteristic1Config GattCharacteristic1 configuration
type GattCharacteristic1Config struct {
	objectPath dbus.ObjectPath
	service    *GattService1
	ID         int
}

// GattCharacteristic1 client
type GattCharacteristic1 struct {
	config      *GattCharacteristic1Config
	properties  *profile.GattCharacteristic1Properties
	descriptors map[dbus.ObjectPath]*GattDescriptor1
	descIndex   int
	notifying   bool
}

//Path return the object path
func (s *GattCharacteristic1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Iface return the Dbus interface
func (s *GattCharacteristic1) Iface() string {
	return bluez.GattCharacteristic1Interface
}

//Properties return the properties of the service
func (s *GattCharacteristic1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	s.properties.Descriptors = s.GetDescriptorPaths()
	p[bluez.GattCharacteristic1Interface] = s.properties
	return p
}

//GetDescriptors return the characteristics of the service
func (s *GattCharacteristic1) GetDescriptors() map[dbus.ObjectPath]*GattDescriptor1 {
	return s.descriptors
}

//GetDescriptorPaths return the characteristics object paths
func (s *GattCharacteristic1) GetDescriptorPaths() []dbus.ObjectPath {
	paths := make([]dbus.ObjectPath, 0)
	for path := range s.descriptors {
		paths = append(paths, path)
	}
	return paths
}

//CreateDescriptor create a new characteristic
func (s *GattCharacteristic1) CreateDescriptor(props *profile.GattDescriptor1Properties) *GattDescriptor1 {
	s.descIndex++
	path := string(s.config.objectPath) + "/desc" + strconv.Itoa(s.descIndex)
	config := &GattDescriptor1Config{
		ID:         s.descIndex,
		objectPath: dbus.ObjectPath(path),
	}
	char := NewGattDescriptor1(config, props)
	return char
}

//AddDescriptor add a characteristic
func (s *GattCharacteristic1) AddDescriptor(char *GattDescriptor1) error {
	s.descriptors[char.Path()] = char
	om := s.config.service.GetApp().GetObjectManager()
	return om.AddObject(char.Path(), char.Properties())
}

//RemoveDescriptor remove a characteristic
func (s *GattCharacteristic1) RemoveDescriptor(char *GattDescriptor1) error {
	if _, ok := s.descriptors[char.Path()]; ok {
		delete(s.descriptors, char.Path())
		om := s.config.service.GetApp().GetObjectManager()
		return om.RemoveObject(char.Path())
	}
	return nil
}

//ReadValue read a value
func (s *GattCharacteristic1) ReadValue(options map[string]interface{}) []byte {
	log.Debug("Characteristic.ReadValue")
	b := make([]byte, 0)
	return b
}

//WriteValue write a value
func (s *GattCharacteristic1) WriteValue(value []byte, options map[string]interface{}) {
	log.Debug("Characteristic.WriteValue")
}

//StartNotify start notification
func (s *GattCharacteristic1) StartNotify() error {
	log.Debug("Characteristic.StartNotify")
	s.notifying = true
	return nil
}

//StopNotify stop notification
func (s *GattCharacteristic1) StopNotify() error {
	log.Debug("Characteristic.StopNotify")
	s.notifying = false
	return nil
}
