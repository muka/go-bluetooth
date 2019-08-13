package service

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

type CharReadCallback func(c *Char, options map[string]interface{}) ([]byte, error)
type CharWriteCallback func(c *Char, value []byte) ([]byte, error)

type Char struct {
	app   *App
	path  dbus.ObjectPath
	descr map[dbus.ObjectPath]*Descr

	props  *gatt.GattCharacteristic1Properties
	iprops *DBusProperties

	readCallback  CharReadCallback
	writeCallback CharWriteCallback
}

func (s *Char) Path() dbus.ObjectPath {
	return s.path
}

func (s *Char) DBusProperties() *DBusProperties {
	return s.iprops
}

func (s *Char) Interface() string {
	return gatt.GattCharacteristic1Interface
}

func (s *Char) Properties() bluez.Properties {
	return s.props
}

func (c *Char) GetDescr() map[dbus.ObjectPath]*Descr {
	return c.descr
}

func (c *Char) App() *App {
	return c.app
}

func (s *Char) RemoveDescr(descr *Descr) error {

	if _, ok := s.descr[descr.Path()]; !ok {
		return nil
	}

	err := descr.Remove()
	if err != nil {
		return err
	}

	delete(s.descr, descr.Path())

	return nil
}

// Expose char to dbus
func (s *Char) Expose() error {
	return ExposeDBusService(s)
}

// Remove char from dbus
func (s *Char) Remove() error {
	return RemoveDBusService(s)
}

// Init new descr
func (s *Char) NewDescr(uuid string) (*Descr, error) {

	descr := new(Descr)
	descr.app = s.App()
	descr.props = NewGattDescriptor1Properties(uuid)
	descr.path = dbus.ObjectPath(
		fmt.Sprintf("%s/descr%d", s.Path(), len(s.GetDescr())),
	)
	iprops, err := NewDBusProperties()
	if err != nil {
		return nil, err
	}
	descr.iprops = iprops

	return descr, nil
}

// Add descr to dbus
func (s *Char) AddDescr(descr *Descr) error {

	err := ExposeDBusService(descr)
	if err != nil {
		return err
	}

	s.descr[descr.Path()] = descr

	return nil
}

// Set the Read callback, called when a client attempt to read
func (s *Char) OnRead(fx CharReadCallback) *Char {
	s.readCallback = fx
	return s
}

// Set the Write callback, called when a client attempt to write
func (s *Char) OnWrite(fx CharWriteCallback) *Char {
	s.writeCallback = fx
	return s
}
