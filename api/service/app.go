package service

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/google/uuid"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

var AppPath = "/org/bluez/%s/app%d"

var UseRandomUUID = false

var baseUUID = "00000000-000%d-1000-8000-00805F9B34FB"

var appCounter = 0

func RandomUUID() (string, error) {
	bUUID := fmt.Sprintf(baseUUID, appCounter)
	UUID, err := uuid.Parse(bUUID)
	if UseRandomUUID {
		UUID, err = uuid.NewRandom()
	}
	return strings.ToUpper(UUID.String()), err
}

// Initialize a new bluetooth service (app)
func NewApp(adapterID string) (*App, error) {

	app := new(App)
	app.adapterID = adapterID
	app.services = make(map[dbus.ObjectPath]*Service)
	app.path = dbus.ObjectPath(
		fmt.Sprintf(
			AppPath,
			app.adapterID,
			appCounter,
		),
	)
	app.advertisement = &advertising.LEAdvertisement1Properties{
		Type: advertising.AdvertisementTypePeripheral,
	}

	appCounter++

	return app, app.init()
}

// Wrap a bluetooth application exposing services
type App struct {
	path          dbus.ObjectPath
	baseUUID      string
	adapterID     string
	conn          *dbus.Conn
	agent         agent.Agent1Client
	objectManager *api.DBusObjectManager
	adapter       *adapter.Adapter1
	services      map[dbus.ObjectPath]*Service
	advertisement *advertising.LEAdvertisement1Properties
	gm            *gatt.GattManager1
}

// return the app dbus path
func (app *App) Path() dbus.ObjectPath {
	return app.path
}

// return the dbus connection
func (app *App) DBusConn() *dbus.Conn {
	return app.conn
}

func (app *App) DBusObjectManager() *api.DBusObjectManager {
	return app.objectManager
}

func (app *App) SetName(name string) {
	app.advertisement.LocalName = name
}

func (app *App) init() error {

	log.Tracef("Exposing %s", app.Path())

	log.Trace("Load adapter")
	a, err := adapter.NewAdapter1FromAdapterID(app.adapterID)
	if err != nil {
		return err
	}
	app.adapter = a

	log.Trace("Register agent")
	agent1, err := app.createAgent()
	if err != nil {
		return err
	}
	app.agent = agent1

	log.Trace("Connecting to DBus")
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	app.conn = conn

	log.Trace("Create object manager")
	om, err := api.NewDBusObjectManager(app.DBusConn())
	if err != nil {
		return err
	}
	app.objectManager = om

	return nil
}

func (app *App) Run() (err error) {

	conn := app.DBusConn()

	_, err = conn.RequestName(
		"org.bluez",
		dbus.NameFlagDoNotQueue&dbus.NameFlagReplaceExisting,
	)
	if err != nil {
		return err
	}

	log.Trace("Exposing Object Manager")
	err = conn.Export(app.objectManager, app.Path(), bluez.ObjectManagerInterface)
	if err != nil {
		return err
	}

	log.Trace("Exposing tree")
	err = app.ExportTree()
	if err != nil {
		return err
	}

	gm, err := gatt.NewGattManager1FromAdapterID(app.adapterID)
	if err != nil {
		return err
	}
	app.gm = gm

	options := map[string]interface{}{}
	err = gm.RegisterApplication(app.Path(), options)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Close() {
	if app.agent != nil {
		err := app.agent.Release()
		if err != nil {
			log.Warnf("Agent1.Release: %s", err)
		}
	}
	if app.gm != nil {
		err1 := app.gm.UnregisterApplication(app.Path())
		if err1 != nil {
			log.Warnf("GattManager1.UnregisterApplication: %s", err1)
		}
	}
}

func (app *App) ExportTree() error {

	childrenNode := make([]introspect.Node, 0)

	for servicePath, service := range app.GetServices() {
		childrenNode = append(childrenNode, introspect.Node{
			Name: string(servicePath)[1:],
		})
		for charPath, char := range service.GetChars() {
			childrenNode = append(childrenNode, introspect.Node{
				Name: string(charPath)[1:],
			})
			for descPath := range char.GetDescr() {
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

	err := app.conn.Export(
		introspect.NewIntrospectable(node),
		app.Path(),
		"org.freedesktop.DBus.Introspectable")

	return err
}
