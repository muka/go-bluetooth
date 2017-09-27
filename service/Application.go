package service

import (
	"errors"
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

	// props, err := NewProperties(config.conn)
	// if err != nil {
	// 	return nil, err
	// }

	s := &Application{
		config:        config,
		objectManager: om,
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
	services      map[dbus.ObjectPath]*GattService1
}

//GetObjectManager return the object manager interface handler
func (app *Application) GetObjectManager() *ObjectManager {
	return app.objectManager
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
	path := string(app.Path()) + "/service" + strconv.Itoa(app.config.serviceIndex)
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

	err = app.exportTree()
	if err != nil {
		return err
	}

	log.Debug("Exposing service to ObjectManager")
	err = app.GetObjectManager().AddObject(service.Path(), service.Properties())
	if err != nil {
		return err
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

		err = app.exportTree()
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

//expose dbus interfaces
func (app *Application) expose() error {

	log.Debugf("Exposing object %s", app.Name())

	conn := app.config.conn
	reply, err := conn.RequestName(app.Name(), dbus.NameFlagReplaceExisting)
	if err != nil {
		log.Debugf("Error requesting object name: %s", err.Error())
		return err
	}

	replym := ""
	switch reply {
	case dbus.RequestNameReplyAlreadyOwner:
		replym = "RequestNameReplyAlreadyOwner"
		break
	case dbus.RequestNameReplyPrimaryOwner:
		replym = "RequestNameReplyPrimaryOwner"
		break
	case dbus.RequestNameReplyExists:
		replym = "RequestNameReplyExists"
		break
	case dbus.RequestNameReplyInQueue:
		replym = "RequestNameReplyInQueue"
		break
	}
	log.Debugf("Name registration reply (%d) %s", reply, replym)

	log.Debugf("Exposing path %s", app.Path())

	// / path
	err = conn.Export(app.objectManager, "/", bluez.ObjectManagerInterface)
	if err != nil {
		return err
	}

	err = app.exportTree()
	if err != nil {
		return err
	}

	dbg("Listening on %s %s", app.Name(), app.Path())

	return nil
}

func (app *Application) exportTree() error {

	childrenNode := make([]introspect.Node, 0)

	for servicePath, service := range app.GetServices() {
		childrenNode = append(childrenNode, introspect.Node{
			Name: string(servicePath)[1:],
		})
		for charPath, char := range service.GetCharacteristics() {
			childrenNode = append(childrenNode, introspect.Node{
				Name: string(charPath)[1:],
			})
			for descPath := range char.GetDescriptors() {
				childrenNode = append(childrenNode, introspect.Node{
					Name: string(descPath)[1:],
				})
			}
		}
	}

	// log.Debugf("child %v", childrenNode)

	// must include also child nodes
	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//ObjectManager
			bluez.ObjectManagerIntrospectData,
		},
		Children: childrenNode,
	}

	err := app.config.conn.Export(
		introspect.NewIntrospectable(node),
		app.Path(),
		"org.freedesktop.DBus.Introspectable")

	return err
}

//Run start the application
func (app *Application) Run() error {

	err := app.expose()
	if err != nil {
		return err
	}

	return nil
}
