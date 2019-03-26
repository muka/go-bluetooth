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

//GattWriteCallback A callback we can register to handle write requests
type GattWriteCallback func(app *Application, service_uuid string, charUUID string, value []byte) error

//GattDescriptorWriteCallback A callback we can register to handle descriptor write requests
type GattDescriptorWriteCallback func(app *Application, service_uuid string, charUUID string, descUUID string, value []byte) error

//GattReadCallback A callback we can register to handle read requests
type GattReadCallback func(app *Application, service_uuid string, charUUID string) ([]byte, error)

//GattDescriptorReadCallback A callback we can register to handle descriptor ead requests
type GattDescriptorReadCallback func(app *Application, service_uuid string, charUUID string, descUUID string) ([]byte, error)

// ApplicationConfig configuration for the bluetooth service
type ApplicationConfig struct {
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

	adMgr         *profile.LEAdvertisingManager1
	advertisement *LEAdvertisement1
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
func (app *Application) CreateService(props *profile.GattService1Properties, advertisedOptional ...bool) (*GattService1, error) {
	app.config.serviceIndex++
	appPath := string(app.Path())
	if appPath == "/" {
		appPath = ""
	}

	advertise := false
	if len(advertisedOptional) > 0 {
		advertise = advertisedOptional[0]
	}

	path := appPath + "/service" + strconv.Itoa(app.config.serviceIndex)
	c := &GattService1Config{
		app:        app,
		objectPath: dbus.ObjectPath(path),
		ID:         app.config.serviceIndex,
		conn:       app.config.conn,
		advertised: advertise,
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

// CallbackError error from a callback
type CallbackError struct {
	msg  string
	code int
}

func (e *CallbackError) Error() string {
	return e.msg
}

//NewCallbackError create a new callback error
func NewCallbackError(code int, msg string) *CallbackError {
	result := &CallbackError{msg: msg, code: code}
	return result
}

//CallbackNotRegistered callback not registered
const CallbackNotRegistered = -1

//CallbackFunctionError callback reported an error
const CallbackFunctionError = -2

//HandleRead Handle application read
func (app *Application) HandleRead(srvUUID string, uuid string) ([]byte, *CallbackError) {
	if app.config.ReadFunc == nil {
		b := make([]byte, 0)
		return b, NewCallbackError(-1, "No callback registered.")
	}

	var cberr *CallbackError
	b, err := app.config.ReadFunc(app, srvUUID, uuid)
	if err != nil {
		cberr = NewCallbackError(-2, err.Error())
	}

	return b, cberr
}

// HandleWrite handle application write
func (app *Application) HandleWrite(srvUUID string, uuid string, value []byte) *CallbackError {
	if app.config.WriteFunc == nil {
		return NewCallbackError(-1, "No callback registered.")
	}

	err := app.config.WriteFunc(app, srvUUID, uuid, value)
	if err != nil {
		return NewCallbackError(-2, err.Error())
	}

	return nil
}

//HandleDescriptorRead handle descriptor read
func (app *Application) HandleDescriptorRead(srvUUID string, charUUID string, descUUID string) ([]byte, *CallbackError) {
	if app.config.DescReadFunc == nil {
		b := make([]byte, 0)
		return b, NewCallbackError(-1, "No callback registered.")
	}

	var cberr *CallbackError
	b, err := app.config.DescReadFunc(app, srvUUID, charUUID, descUUID)
	if err != nil {
		cberr = NewCallbackError(-2, err.Error())
	}

	return b, cberr
}

//HandleDescriptorWrite handle descriptor write
func (app *Application) HandleDescriptorWrite(srvUUID string, charUUID string, descUUID string, value []byte) *CallbackError {
	if app.config.DescWriteFunc == nil {
		return NewCallbackError(-1, "No callback registered.")
	}

	err := app.config.DescWriteFunc(app, srvUUID, charUUID, descUUID, value)
	if err != nil {
		return NewCallbackError(-2, err.Error())
	}

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

//StartAdvertising advertise information for a service
func (app *Application) StartAdvertising(deviceInterface string, appID string) error {
	if app.advertisement != nil && app.adMgr != nil {
		// Already advertising
		return nil
	}

	path := "/org/bluez/advertisement/" + appID

	config := &LEAdvertisement1Config{
		conn:       app.config.conn,
		objectPath: dbus.ObjectPath(path),
	}

	serviceUUIDs := make([]string, 0)

	for _, serv := range app.services {
		if serv.Advertised() {
			serviceUUIDs = append(serviceUUIDs, serv.properties.UUID)
		}
	}

	props := &profile.LEAdvertisement1Properties{
		Type:         "peripheral",
		LocalName:    app.config.LocalName,
		ServiceUUIDs: serviceUUIDs,
	}

	var err error

	app.advertisement, err = NewLEAdvertisement1(config, props)
	if err != nil {
		app.advertisement = nil
		return err
	}

	err = app.advertisement.Expose()
	if err != nil {
		app.advertisement = nil
		return err
	}

	options := make(map[string]interface{})

	app.adMgr = profile.NewLEAdvertisingManager1(deviceInterface)

	err = app.adMgr.RegisterAdvertisement(path, options)
	if err != nil {
		app.advertisement = nil
		app.adMgr = nil
		return err
	}

	adapter := profile.NewAdapter1(deviceInterface)
	err = adapter.SetProperty("Discoverable", dbus.MakeVariant(true))
	if err != nil {
		return err
	}

	err = adapter.SetProperty("Powered", dbus.MakeVariant(true))
	if err != nil {
		return err
	}

	return nil
}

//StopAdvertising stop advertising information on a service
func (app *Application) StopAdvertising() error {
	if app.advertisement == nil || app.adMgr == nil {
		// Not advertising
		return nil
	}

	err := app.adMgr.UnregisterAdvertisement(string(app.advertisement.config.objectPath))

	app.advertisement = nil
	app.adMgr = nil

	if err != nil {
		return err
	}

	return nil
}
