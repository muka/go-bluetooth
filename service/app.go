package service

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
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

	a, err := adapter.GetAdapter(config.AdapterID)
	if err != nil {
		return nil, err
	}

	s := &Application{
		config:        config,
		objectManager: om,
		adapter:       a,
		services:      make(map[dbus.ObjectPath]*GattService1),
	}

	return s, nil
}

// ApplicationConfig configuration for the bluetooth service
type ApplicationConfig struct {
	AdapterID    string
	UUIDSuffix   string
	UUID         string
	conn         *dbus.Conn
	ObjectName   string
	ObjectPath   dbus.ObjectPath
	serviceIndex int
	LocalName    string

	WriteFunc     GattWriteCallback
	ReadFunc      GattReadCallback
	DescWriteFunc GattDescriptorWriteCallback
	DescReadFunc  GattDescriptorReadCallback
}

// Application a bluetooth service exposed by bluez
type Application struct {
	config        *ApplicationConfig
	objectManager *ObjectManager
	services      map[dbus.ObjectPath]*GattService1
	adapter       *adapter.Adapter1
	adMgr         *advertising.LEAdvertisingManager1
	advertisement *LEAdvertisement1
}

//GetAdapter return the Adapter1 interface instance
func (app *Application) GetAdapter() *adapter.Adapter1 {
	return app.adapter
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

//Run start the application
func (app *Application) Run() error {

	err := app.expose()
	if err != nil {
		return err
	}

	return nil
}
