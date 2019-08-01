package service

import (
	"strconv"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

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
func (s *GattService1) CreateCharacteristic(props *gatt.GattCharacteristic1Properties) (*GattCharacteristic1, error) {
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

	_, ok := s.characteristics[char.Path()]
	s.characteristics[char.Path()] = char
	if !ok {
		s.properties.Characteristics = append(s.properties.Characteristics, char.Path())
	}

	err := char.Expose()
	if err != nil {
		return err
	}

	err = s.GetApp().exportTree()
	if err != nil {
		return err
	}

	om := s.config.app.GetObjectManager()
	return om.AddObject(char.Path(), map[string]bluez.Properties{
		char.Interface(): char.Properties(),
	})
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
