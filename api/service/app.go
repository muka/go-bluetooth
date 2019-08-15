package service

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/google/uuid"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	log "github.com/sirupsen/logrus"
)

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
			"/org/bluez/%s/app%d",
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
	objectManager *ObjectManager
	adapter       *adapter.Adapter1
	services      map[dbus.ObjectPath]*Service
	advertisement *advertising.LEAdvertisement1Properties
}

// return the app dbus path
func (app *App) Path() dbus.ObjectPath {
	return app.path
}

// return the dbus connection
func (app *App) DBusConn() *dbus.Conn {
	return app.conn
}

func (app *App) ObjectManager() *ObjectManager {
	return app.objectManager
}

func (app *App) SetName(name string) {
	app.advertisement.LocalName = name
}

// Start the app
func (app *App) init() error {

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

	log.Trace("Register object manager")
	om, err := NewObjectManager(app.DBusConn())
	if err != nil {
		return err
	}
	app.objectManager = om

	_, err = conn.RequestName(
		"org.bluez",
		dbus.NameFlagDoNotQueue&dbus.NameFlagReplaceExisting,
	)
	if err != nil {
		return err
	}

	log.Tracef("Exposing %s", app.Path())

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

func (app *App) Close() {
	err := app.agent.Release()
	if err != nil {
		log.Warnf("agent1.Release: %s", err)
	}
}

func (app *App) exportTree() error {

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

	// log.Tracef("Exposing XML Tree on %s: \n\n%s", app.Path(), *node)

	err := app.conn.ExportSubtree(
		introspect.NewIntrospectable(node),
		app.Path(),
		"org.freedesktop.DBus.Introspectable")

	return err
}
