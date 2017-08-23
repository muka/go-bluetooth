package service

import (
	"errors"
	"strconv"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/satori/go.uuid"
	"github.com/tj/go-debug"
)

var dbg = debug.Debug("bluetooth:server")

//NewApplication instantiate a new application service
func NewApplication(config *ApplicationConfig) (*Application, error) {

	if config.conn == nil {
		conn, err := dbus.SystemBus()
		if err != nil {
			return nil, err
		}
		config.conn = conn
	}

	om, err := NewObjectManager(config.conn)
	if err != nil {
		return nil, err
	}

	s := &Application{
		config,
		om,
		make(map[dbus.ObjectPath]*GattService1),
	}

	return s, nil
}

// ApplicationConfig configuration for the bluetooth service
type ApplicationConfig struct {
	conn         *dbus.Conn
	objectName   string
	objectPath   dbus.ObjectPath
	serviceIndex int
}

// Application a bluetooth service exposed by bluez
type Application struct {
	config        *ApplicationConfig
	objectManager *ObjectManager
	services      map[dbus.ObjectPath]*GattService1
}

//GetObjectManager return the object manager
func (app *Application) GetObjectManager() *ObjectManager {
	return app.objectManager
}

//expose dbus interfaces
func (app *Application) expose() error {

	dbg("Exposing object %s", app.config.objectName)
	conn := app.config.conn
	reply, err := conn.RequestName(app.config.objectName, dbus.NameFlagDoNotQueue)
	if err != nil {
		dbg("Error requesting object name: %s", err.Error())
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Requested name has been already taken")
	}

	// f := foo("Bar!")
	// conn.Export(f, "/com/github/guelfey/Demo", "com.github.guelfey.Demo")

	intro := ""
	conn.Export(
		introspect.Introspectable(intro),
		app.config.objectPath,
		"org.freedesktop.DBus.Introspectable")

	dbg("Listening on %s %s", app.config.objectName, app.config.objectPath)

	return nil
}

//Path return the object path
func (app *Application) Path() dbus.ObjectPath {
	return app.config.objectPath
}

//Run start the application
func (app *Application) Run() error {
	return nil
}

//AddService add service to expose
func (app *Application) AddService(service *GattService1) error {
	app.services[service.Path()] = service
	err := app.objectManager.AddObject(service.Path(), service.Properties())
	return err
}

//RemoveService remove an exposed service
func (app *Application) RemoveService(service *GattService1) error {

	if _, ok := app.services[service.Path()]; ok {
		delete(app.services, service.Path())
		err := app.objectManager.RemoveObject(service.Path())
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateUUID generate a UUIDv4
func (app *Application) GenerateUUID() string {
	return uuid.NewV4().String()
}

//CreateService create a new GattService1 instance
func (app *Application) CreateService(props *profile.GattService1Properties) (*GattService1, error) {
	path := string(app.config.objectPath) + strconv.Itoa(app.config.serviceIndex)
	c := &GattService1Config{
		app:      app,
		basePath: dbus.ObjectPath(path),
		ID:       app.config.serviceIndex,
	}
	app.config.serviceIndex++
	s := NewGattService1(c, props)
	return s, nil
}
