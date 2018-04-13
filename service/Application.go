package service

import (
	"errors"
	"strconv"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

//UUIDSuffix fixed 128bit UUID [0000]+[xxxx]+[-0000-1000-8000-00805F9B34FB]
const UUIDSuffix = "-0000-1000-8000-00805F9B34FB"

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
	UUIDSuffix   string
	UUID         string
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

// GenerateUUID generate a 128bit UUID
func (app *Application) GenerateUUID(uuidVal string) string {
	base := app.config.UUID
	if len(uuidVal) == 8 {
		base = ""
	}
	return base + uuidVal + app.config.UUIDSuffix
}

//CreateService create a new GattService1 instance
func (app *Application) CreateService(props *profile.GattService1Properties) (*GattService1, error) {
	app.config.serviceIndex++
	appPath := string(app.Path())
	if appPath == "/" {
		appPath = ""
	}
	path := appPath + "/service" + strconv.Itoa(app.config.serviceIndex)
	c := &GattService1Config{
		app:        app,
		objectPath: dbus.ObjectPath(path),
		ID:         app.config.serviceIndex,
		conn:       app.config.conn,
	}
	s, err := NewGattService1(c, props)
	return s, err
}

//AddService add service to expose
func (app *Application) AddService(service *GattService1) error {

	app.services[service.Path()] = service

	err := service.Expose()
	if err != nil {
		return err
	}

	err = app.exportTree()
	if err != nil {
		return err
	}

	err = app.GetObjectManager().AddObject(service.Path(), service.Properties())
	if err != nil {
		return err
	}

	return err
}

//RemoveService remove an exposed service
func (app *Application) RemoveService(service *GattService1) error {
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

	conn := app.config.conn
	_, err := conn.RequestName(app.Name(), dbus.NameFlagDoNotQueue&dbus.NameFlagReplaceExisting)
	if err != nil {
		return err
	}

	// / path
	err = conn.Export(app.objectManager, app.Path(), bluez.ObjectManagerInterface)
	if err != nil {
		return err
	}

	err = app.exportTree()
	if err != nil {
		return err
	}

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
