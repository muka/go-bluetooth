package service

import (
	"errors"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/satori/go.uuid"
	"github.com/tj/go-debug"
)

var dbg = debug.Debug("bluetooth:server")

//NewApplication instantiate a new application service
func NewApplication(config *ApplicationConfig) (*Application, error) {

	if config.ObjectName == "" {
		return nil, errors.New("objectName is required")
	}
	if config.ObjectPath == "" {
		return nil, errors.New("objectPath is required")
	}

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

	props, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	s := &Application{
		config,
		om,
		props,
		make(map[dbus.ObjectPath]*GattService1),
	}

	return s, nil
}

// ApplicationConfig configuration for the bluetooth service
type ApplicationConfig struct {
	conn         *dbus.Conn
	ObjectName   string
	ObjectPath   dbus.ObjectPath
	serviceIndex int
}

// Application a bluetooth service exposed by bluez
type Application struct {
	config        *ApplicationConfig
	objectManager *ObjectManager
	properties    *Properties
	services      map[dbus.ObjectPath]*GattService1
}

//GetObjectManager return the object manager interface handler
func (app *Application) GetObjectManager() *ObjectManager {
	return app.objectManager
}

//GetProperties return the properties interface handler
func (app *Application) GetProperties() *Properties {
	return app.properties
}

//Path return the object path
func (app *Application) Path() dbus.ObjectPath {
	return app.config.ObjectPath
}

//Name return the object name
func (app *Application) Name() string {
	return app.config.ObjectName
}

// GenerateUUID generate a UUIDv4
func (app *Application) GenerateUUID() string {
	return uuid.NewV4().String()
}

//CreateService create a new GattService1 instance
func (app *Application) CreateService(props *profile.GattService1Properties) (*GattService1, error) {
	path := string(app.Path()) + strconv.Itoa(app.config.serviceIndex)
	c := &GattService1Config{
		app:      app,
		basePath: dbus.ObjectPath(path),
		ID:       app.config.serviceIndex,
	}
	app.config.serviceIndex++
	s := NewGattService1(c, props)
	return s, nil
}

//AddService add service to expose
func (app *Application) AddService(service *GattService1) error {
	log.Debugf("Adding service %s", service.Path())
	app.services[service.Path()] = service
	err := app.objectManager.AddObject(service.Path(), service.Properties())
	return err
}

//RemoveService remove an exposed service
func (app *Application) RemoveService(service *GattService1) error {
	log.Debugf("Removing service %s", service.Path())
	if _, ok := app.services[service.Path()]; ok {
		delete(app.services, service.Path())
		//TODO: remove chars + descritptors too
		err := app.objectManager.RemoveObject(service.Path())
		if err != nil {
			return err
		}
	}

	return nil
}

//expose dbus interfaces
func (app *Application) expose() error {

	log.Debugf("Exposing object %s", app.Name())
	conn := app.config.conn
	reply, err := conn.RequestName(app.Name(), dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Debugf("Error requesting object name: %s", err.Error())
		return err
	}

	log.Debugf("Name registration reply %d", reply)
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("Requested name has been already taken (%d)", reply)
	}

	conn.Export(app.objectManager, "/", bluez.ObjectManagerInterface)
	conn.Export(app.properties, "/", bluez.PropertiesInterface)

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Properties
			bluez.PropertiesIntrospectData,
			//ObjectManager
			bluez.ObjectManagerIntrospectData,
			//Introspect
			introspect.IntrospectData,
		},
	}

	log.Debugf("Exposing dbus service\n%s\n", node)

	conn.Export(
		introspect.NewIntrospectable(node),
		app.Path(),
		"org.freedesktop.DBus.Introspectable")

	dbg("Listening on %s %s", app.Name(), app.Path())

	return nil
}

//Run start the application
func (app *Application) Run() error {

	err := app.expose()
	if err != nil {
		return err
	}

	return nil
}
