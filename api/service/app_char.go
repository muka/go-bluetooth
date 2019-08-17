package service

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

type CharReadCallback func(c *Char, options map[string]interface{}) ([]byte, error)
type CharWriteCallback func(c *Char, value []byte) ([]byte, error)

type Char struct {
	ID      int
	app     *App
	service *Service

	path  dbus.ObjectPath
	descr map[dbus.ObjectPath]*Descr

	Properties *gatt.GattCharacteristic1Properties
	iprops     *api.DBusProperties

	readCallback  CharReadCallback
	writeCallback CharWriteCallback
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

// Init new descr
func (s *Char) NewDescr() (*Descr, error) {

	descr := new(Descr)
	descr.ID = s.ID + len(s.descr) + 10

	baseUUID := "%08x" + s.Properties.UUID[8:]
	uuid := fmt.Sprintf(baseUUID, descr.ID)

	descr.app = s.App()
	descr.char = s
	descr.Properties = NewGattDescriptor1Properties(uuid)
	descr.path = dbus.ObjectPath(
		fmt.Sprintf("%s/descr%d", s.Path(), len(s.GetDescr())),
	)
	iprops, err := api.NewDBusProperties(s.App().DBusConn())
	if err != nil {
		return nil, err
	}
	descr.iprops = iprops

	return descr, nil
}

// Add descr to dbus
func (s *Char) AddDescr(descr *Descr) error {

	err := api.ExposeDBusService(descr)
	if err != nil {
		return err
	}

	s.descr[descr.Path()] = descr

	// log.Tracef("Added GATT Descriptor ID=%d %s", descr.ID, descr.Path())

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

// start notification session
func (s *Char) StartNotify() *dbus.Error {
	log.Debug("Char.StartNotify")
	return nil
}

// stop notification session
func (s *Char) StopNotify() *dbus.Error {
	log.Debug("Char.StopNotify")
	return nil
}
