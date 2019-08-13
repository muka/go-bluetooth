package service

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

type Service struct {
	app   *App
	path  dbus.ObjectPath
	props *gatt.GattService1Properties
	chars map[dbus.ObjectPath]*Char
}

func (s *Service) Path() dbus.ObjectPath {
	return s.path
}

func (s *Service) Interface() string {
	return gatt.GattService1Interface
}

func (s *Service) Properties() bluez.Properties {
	return s.props
}

func (s *Service) App() *App {
	return s.app
}

// Expose service to dbus
func (s *Service) Expose() error {
	return ExposeService(s)
}

// Remove service from dbus
func (s *Service) Remove() error {
	return RemoveService(s)
}

func (s *Service) GetChars() map[dbus.ObjectPath]*Char {
	return s.chars
}

// Create a new characteristic
func (s *Service) NewChar(uuid string) *Char {

	char := new(Char)
	char.path = dbus.ObjectPath(
		fmt.Sprintf("%s/char%d", s.Path(), len(s.GetChars())),
	)
	char.app = s.App()
	char.descr = make(map[dbus.ObjectPath]*Descr)
	char.props = NewGattCharacteristic1Properties(uuid)

	return char
}

func (s *Service) AddChar(char *Char) error {

	s.chars[char.Path()] = char

	iprops, err := NewDBusProperties()
	if err != nil {
		return err
	}
	char.iprops = iprops

	err = ExposeService(char)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveChar(char *Char) error {
	if _, ok := s.chars[char.Path()]; !ok {
		return nil
	}

	for _, descr := range char.GetDescr() {
		err := char.RemoveDescr(descr)
		if err != nil {
			return err
		}
	}

	// remove the char from the three
	err := s.app.objectManager.RemoveObject(s.Path())
	if err != nil {
		return err
	}

	err = s.app.exportTree()
	if err != nil {
		return err
	}

	delete(s.chars, char.Path())

	return nil
}
