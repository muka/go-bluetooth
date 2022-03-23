package service

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

type CharReadCallback func(c *Char, options map[string]interface{}) ([]byte, error)
type CharWriteCallback func(c *Char, value []byte) ([]byte, error)
type CharNotifyCallback func(c *Char, notify bool) error

type Char struct {
	UUID    string
	app     *App
	service *Service

	path  dbus.ObjectPath
	descr map[dbus.ObjectPath]*Descr

	Properties *gatt.GattCharacteristic1Properties
	iprops     *api.DBusProperties

	readCallback   CharReadCallback
	writeCallback  CharWriteCallback
	notifyCallback CharNotifyCallback
}

func (s *Char) Path() dbus.ObjectPath {
	return s.path
}

func (s *Char) DBusProperties() *api.DBusProperties {
	return s.iprops
}

func (s *Char) Interface() string {
	return gatt.GattCharacteristic1Interface
}

func (s *Char) GetProperties() bluez.Properties {
	descr := []dbus.ObjectPath{}
	for dpath := range s.descr {
		descr = append(descr, dpath)
	}
	s.Properties.Descriptors = descr
	s.Properties.Service = s.Service().Path()

	return s.Properties
}

func (c *Char) GetDescr() map[dbus.ObjectPath]*Descr {
	return c.descr
}

func (c *Char) App() *App {
	return c.app
}

func (c *Char) Service() *Service {
	return c.service
}

func (s *Char) DBusObjectManager() *api.DBusObjectManager {
	return s.App().DBusObjectManager()
}

func (s *Char) DBusConn() *dbus.Conn {
	return s.App().DBusConn()
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
	return api.ExposeDBusService(s)
}

// Remove char from dbus
func (s *Char) Remove() error {
	return api.RemoveDBusService(s)
}

// NewDescr Init new descr
func (s *Char) NewDescr(uuid string) (*Descr, error) {

	descr := new(Descr)
	descr.UUID = s.App().GenerateUUID(uuid)

	descr.app = s.App()
	descr.char = s
	descr.Properties = NewGattDescriptor1Properties(descr.UUID)
	descr.path = dbus.ObjectPath(
		fmt.Sprintf("%s/descriptor%d", s.Path(), len(s.GetDescr())),
	)
	iprops, err := api.NewDBusProperties(s.App().DBusConn())
	if err != nil {
		return nil, err
	}
	descr.iprops = iprops

	return descr, nil
}

// AddDescr Add descr to dbus
func (s *Char) AddDescr(descr *Descr) error {

	err := api.ExposeDBusService(descr)
	if err != nil {
		return err
	}

	s.descr[descr.Path()] = descr

	err = s.DBusObjectManager().AddObject(descr.Path(), map[string]bluez.Properties{
		descr.Interface(): descr.GetProperties(),
	})
	if err != nil {
		return err
	}

	// update OM char too
	err = s.DBusObjectManager().AddObject(s.Path(), map[string]bluez.Properties{
		s.Interface(): s.GetProperties(),
	})
	if err != nil {
		return err
	}

	log.Tracef("Added GATT Descriptor UUID=%s %s", descr.UUID, descr.Path())

	err = s.App().ExportTree()
	return err
}

// OnRead Set the Read callback, called when a client attempt to read
func (s *Char) OnRead(fx CharReadCallback) *Char {
	s.readCallback = fx
	return s
}

// OnWrite Set the Write callback, called when a client attempt to write
func (s *Char) OnWrite(fx CharWriteCallback) *Char {
	s.writeCallback = fx
	return s
}

// OnNotify Set the Notify callback, called when a client attempt to start/stop notifications for a characteristic
func (s *Char) OnNotify(fx CharNotifyCallback) *Char {
	s.notifyCallback = fx
	return s
}
