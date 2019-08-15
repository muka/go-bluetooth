package service

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	ID         int
	app        *App
	path       dbus.ObjectPath
	Properties *gatt.GattService1Properties
	chars      map[dbus.ObjectPath]*Char
	iprops     *DBusProperties
}

func (s *Service) DBusProperties() *DBusProperties {
	return s.iprops
}

func (s *Service) Path() dbus.ObjectPath {
	return s.path
}

func (s *Service) Interface() string {
	return gatt.GattService1Interface
}

func (s *Service) GetProperties() bluez.Properties {
	return s.Properties
}

func (s *Service) App() *App {
	return s.app
}

// Expose service to dbus
func (s *Service) Expose() error {
	return ExposeDBusService(s)
}

// Remove service from dbus
func (s *Service) Remove() error {
	return RemoveDBusService(s)
}

func (s *Service) GetChars() map[dbus.ObjectPath]*Char {
	return s.chars
}

// Create a new characteristic
func (s *Service) NewChar() (*Char, error) {

	char := new(Char)
	char.ID = s.ID + len(s.chars) + 100

	serviceUUID := "%08x" + s.Properties.UUID[8:]
	uuid := fmt.Sprintf(serviceUUID, char.ID)

	char.path = dbus.ObjectPath(
		fmt.Sprintf("%s/char%d", s.Path(), len(s.GetChars())),
	)
	char.app = s.App()
	char.service = s
	char.descr = make(map[dbus.ObjectPath]*Descr)
	char.Properties = NewGattCharacteristic1Properties(uuid)

	iprops, err := NewDBusProperties(s.App().DBusConn())
	if err != nil {
		return nil, err
	}
	char.iprops = iprops

	return char, nil
}

func (s *Service) AddChar(char *Char) error {

	s.chars[char.Path()] = char

	err := ExposeDBusService(char)
	if err != nil {
		return err
	}

	log.Tracef("Added GATT Characteristic ID=%d %s", char.ID, char.Properties.UUID)

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
