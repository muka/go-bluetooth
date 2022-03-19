package service

import (
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
)

func (app *App) AdapterID() string {
	return app.adapterID
}

func (app *App) Adapter() *adapter.Adapter1 {
	return app.adapter
}

func (app *App) Agent() agent.Agent1Client {
	return app.agent
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

func (app *App) SetLEAdvertisement(adv *advertising.LEAdvertisement1Properties) {
	app.advertisement = adv
}
