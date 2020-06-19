package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

// AppPath default app path
var AppPath = "/%s/apps/%d"

var appCounter = 0

// AppOptions contains App options
type AppOptions struct {
	AdapterID         string
	AgentCaps         string
	AgentSetAsDefault bool
	UUIDSuffix        string
	UUID              string
}

// NewApp initialize a new bluetooth service (app)
func NewApp(options AppOptions) (*App, error) {

	app := new(App)
	if options.AdapterID == "" {
		return nil, errors.New("options.AdapterID is required")
	}

	app.Options = options

	if app.Options.UUIDSuffix == "" {
		app.Options.UUIDSuffix = "-0000-1000-8000-00805F9B34FB"
	}
	if app.Options.UUID == "" {
		app.Options.UUID = "1234"
	}

	app.adapterID = app.Options.AdapterID
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

	if app.Options.AgentCaps == "" {
		app.Options.AgentCaps = agent.CapKeyboardDisplay
	}

	appCounter++

	return app, app.init()
}

// App wraps a bluetooth application exposing services
type App struct {
	path    dbus.ObjectPath
	Options AppOptions

	adapterID string
	adapter   *adapter.Adapter1

	agent agent.Agent1Client

	conn          *dbus.Conn
	objectManager *api.DBusObjectManager
	services      map[dbus.ObjectPath]*Service
	advertisement *advertising.LEAdvertisement1Properties
	gm            *gatt.GattManager1
}

func (app *App) init() error {

	// log.Tracef("Exposing %s", app.Path())

	// log.Trace("Load adapter")
	a, err := adapter.NewAdapter1FromAdapterID(app.adapterID)
	if err != nil {
		return err
	}
	app.adapter = a

	agent1, err := app.createAgent()
	if err != nil {
		return err
	}
	app.agent = agent1

	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	app.conn = conn

	om, err := api.NewDBusObjectManager(app.DBusConn())
	if err != nil {
		return err
	}
	app.objectManager = om

	return err
}

// GenerateUUID generate a 128bit UUID
func (app *App) GenerateUUID(uuidVal string) string {
	base := app.Options.UUID
	if len(uuidVal) == 8 {
		base = ""
	}
	return base + uuidVal + app.Options.UUIDSuffix
}

// GetAdapter return the adapter in use
func (app *App) GetAdapter() *adapter.Adapter1 {
	return app.adapter
}

// Expose children services, chars and descriptors
func (app *App) extractChildren() (children []introspect.Node) {
	for _, service := range app.GetServices() {
		childPath := strings.ReplaceAll(string(service.Path()), string(app.Path())+"/", "")
		children = append(children, introspect.Node{
			Name: childPath,
		})
		// chars
		for _, char := range service.GetChars() {
			childPath := strings.ReplaceAll(string(char.Path()), string(app.Path())+"/", "")
			children = append(children, introspect.Node{
				Name: childPath,
			})
			// descrs
			for _, descr := range char.GetDescr() {
				childPath := strings.ReplaceAll(string(descr.Path()), string(app.Path())+"/", "")
				children = append(children, introspect.Node{
					Name: childPath,
				})
			}
		}
	}
	return children
}

// ExportTree update introspection data
func (app *App) ExportTree() (err error) {

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//ObjectManager
			bluez.ObjectManagerIntrospectData,
		},
		Children: app.extractChildren(),
	}

	introspectable := introspect.NewIntrospectable(node)
	err = app.conn.Export(
		introspectable,
		app.Path(),
		"org.freedesktop.DBus.Introspectable",
	)

	return err
}

// Run initialize the application
func (app *App) Run() (err error) {

	log.Tracef("Expose %s (%s)", app.Path(), bluez.ObjectManagerInterface)
	err = app.conn.Export(app.objectManager, app.Path(), bluez.ObjectManagerInterface)
	if err != nil {
		return err
	}

	err = app.ExportTree()
	if err != nil {
		return err
	}

	err = app.ExposeAgent(app.Options.AgentCaps, app.Options.AgentSetAsDefault)
	if err != nil {
		return fmt.Errorf("ExposeAgent: %s", err)
	}

	gm, err := gatt.NewGattManager1FromAdapterID(app.adapterID)
	if err != nil {
		return err
	}
	app.gm = gm

	options := map[string]interface{}{}
	err = gm.RegisterApplication(app.Path(), options)

	return err
}

// Close close the app
func (app *App) Close() {

	if app.agent != nil {

		err := agent.RemoveAgent(app.agent)
		if err != nil {
			log.Warnf("RemoveAgent: %s", err)
		}

		// err =
		app.agent.Release()
		// if err != nil {
		// 	log.Warnf("Agent1.Release: %s", err)
		// }
	}

	if app.gm != nil {
		err1 := app.gm.UnregisterApplication(app.Path())
		if err1 != nil {
			log.Warnf("GattManager1.UnregisterApplication: %s", err1)
		}
	}
}
