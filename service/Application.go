package service

import (
	"errors"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
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
		config:        config,
		objectManager: om,
		properties:    props,
		services:      make(map[dbus.ObjectPath]*GattService1),
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
	app.config.serviceIndex++
	path := string(app.Path()) + "service" + strconv.Itoa(app.config.serviceIndex)
	c := &GattService1Config{
		app:        app,
		objectPath: dbus.ObjectPath(path),
		ID:         app.config.serviceIndex,
		conn:       app.config.conn,
	}
	s, err := NewGattService1(c, props)
	log.Debugf("Created service %s", path)
	return s, err
}

//AddService add service to expose
func (app *Application) AddService(service *GattService1) error {

	log.Debugf("Adding service %s", service.Path())
	app.services[service.Path()] = service

	err := service.Expose()
	if err != nil {
		return err
	}

	err = app.GetObjectManager().AddObject(service.Path(), service.Properties())
	if err != nil {
		return err
	}

	for iface, props := range service.Properties() {
		app.GetProperties().AddProperties(iface, props)
	}

	return err
}

//RemoveService remove an exposed service
func (app *Application) RemoveService(service *GattService1) error {
	log.Debugf("Removing service %s", service.Path())
	if _, ok := app.services[service.Path()]; ok {
		delete(app.services, service.Path())
		err := app.GetObjectManager().RemoveObject(service.Path())
		//TODO: remove chars + descritptors too
		if err != nil {
			return err
		}
	}
	return nil
}

//GetServices return the registered services
func (app *Application) GetServices() map[dbus.ObjectPath]*GattService1 {
	return app.services
}

func getNode(idata []introspect.Interface) *introspect.Node {
	rootNode := &introspect.Node{
		Interfaces: append([]introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//Properties
			prop.IntrospectData,
		}, idata...),
	}
	return rootNode
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

	log.Debugf("Exposing path %s", app.Path())

	// / path
	err = conn.Export(app.objectManager, "/", bluez.ObjectManagerInterface)
	if err != nil {
		return err
	}
	err = conn.Export(app.properties, "/", bluez.PropertiesInterface)
	if err != nil {
		return err
	}

	rootNode := getNode([]introspect.Interface{
		//ObjectManager
		bluez.ObjectManagerIntrospectData,
	})

	err = conn.Export(
		introspect.NewIntrospectable(rootNode),
		app.Path(),
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	dbg("Listening on %s %s", app.Name(), app.Path())

	return nil
}

//Run start the application
func (app *Application) Run() error {

	err := app.expose()
	if err != nil {
		return err
	}

	app.properties.Expose(app.Path())

	return nil
}
