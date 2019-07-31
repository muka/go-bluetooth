package service

import (
	"strconv"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

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
func (s *GattCharacteristic1) CreateDescriptor(props *gatt.GattDescriptor1Properties) (*GattDescriptor1, error) {
	s.descIndex++
	path := string(s.config.objectPath) + "/desc" + strconv.Itoa(s.descIndex)
	config := &GattDescriptor1Config{
		ID:             s.descIndex,
		objectPath:     dbus.ObjectPath(path),
		conn:           s.config.conn,
		characteristic: s,
	}

	props.Characteristic = config.objectPath

	desc, err := NewGattDescriptor1(config, props)
	return desc, err
}

//AddDescriptor add a characteristic
func (s *GattCharacteristic1) AddDescriptor(desc *GattDescriptor1) error {

	s.descriptors[desc.Path()] = desc

	err := desc.Expose()
	if err != nil {
		return err
	}

	err = s.config.service.GetApp().exportTree()
	if err != nil {
		return err
	}

	om := s.config.service.GetApp().GetObjectManager()
	return om.AddObject(desc.Path(), map[string]bluez.Properties{
		desc.Interface(): desc.Properties(),
	})
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
