package service

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	UUID       string
	app        *App
	path       dbus.ObjectPath
	Properties *gatt.GattService1Properties
	chars      map[dbus.ObjectPath]*Char
	iprops     *api.DBusProperties
}

func (s *Service) DBusProperties() *api.DBusProperties {
	return s.iprops
}

func (s *Service) Path() dbus.ObjectPath {
	return s.path
}

func (s *Service) Interface() string {
	return gatt.GattService1Interface
}

func (s *Service) GetProperties() bluez.Properties {

	chars := []dbus.ObjectPath{}
	for cpath := range s.chars {
		chars = append(chars, cpath)
	}
	s.Properties.Characteristics = chars

	return s.Properties
}

func (s *Service) App() *App {
	return s.app
}

func (s *Service) DBusObjectManager() *api.DBusObjectManager {
	return s.App().DBusObjectManager()
}

func (s *Service) DBusConn() *dbus.Conn {
	return s.App().DBusConn()
}

// Expose service to dbus
func (s *Service) Expose() error {
	return api.ExposeDBusService(s)
}

// Remove service from dbus
func (s *Service) Remove() error {
	return api.RemoveDBusService(s)
}

func (s *Service) GetChars() map[dbus.ObjectPath]*Char {
	return s.chars
}

// NewChar Create a new characteristic
func (s *Service) NewChar(uuid string) (*Char, error) {

	char := new(Char)
	char.UUID = s.App().GenerateUUID(uuid)

	char.path = dbus.ObjectPath(
		fmt.Sprintf("%s/char%d", s.Path(), len(s.GetChars())),
	)
	char.app = s.App()
	char.service = s
	char.descr = make(map[dbus.ObjectPath]*Descr)
	char.Properties = NewGattCharacteristic1Properties(char.UUID)

	iprops, err := api.NewDBusProperties(s.App().DBusConn())
	if err != nil {
		return nil, err
	}
	char.iprops = iprops

	return char, nil
}

func (s *Service) AddChar(char *Char) error {

	s.chars[char.Path()] = char

	err := api.ExposeDBusService(char)
	if err != nil {
		return err
	}

	err = s.DBusObjectManager().AddObject(char.Path(), map[string]bluez.Properties{
		char.Interface(): char.GetProperties(),
	})
	if err != nil {
		return err
	}

	// update OM service rapresentation also
	err = s.DBusObjectManager().AddObject(s.Path(), map[string]bluez.Properties{
		s.Interface(): s.GetProperties(),
	})
	if err != nil {
		return err
	}

	log.Tracef("Added GATT Characteristic UUID=%s %s", char.UUID, char.Path())

	err = s.App().ExportTree()
	return err
}

// RemoveChar remove a characteristic from the service
func (s *Service) RemoveChar(char *Char) error {
	// todo unregister properties
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
	err := s.DBusObjectManager().RemoveObject(s.Path())
	if err != nil {
		return err
	}

	// update OM service rapresentation also
	err = s.DBusObjectManager().AddObject(s.Path(), map[string]bluez.Properties{
		s.Interface(): s.GetProperties(),
	})
	if err != nil {
		return err
	}

	// err = s.ExportTree()
	// if err != nil {
	// 	return err
	// }

	delete(s.chars, char.Path())

	return nil
}
