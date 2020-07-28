package api

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	log "github.com/sirupsen/logrus"
)

// const baseAdvertismentPath = "/org/bluez/%s/apps/advertisement%d"
const BaseAdvertismentPath = "/go_bluetooth/%s/advertisement/%d"

var advertisingCount int = -1

func nextAdvertismentPath(adapterID string) dbus.ObjectPath {
	advertisingCount++
	return dbus.ObjectPath(fmt.Sprintf(BaseAdvertismentPath, adapterID, advertisingCount))
}

func decreaseAdvertismentCounter() {
	advertisingCount--
	if advertisingCount < -1 {
		advertisingCount = -1
	}
}

type Advertisement struct {
	path          dbus.ObjectPath
	objectManager *DBusObjectManager
	iprops        *DBusProperties
	conn          *dbus.Conn
	props         *advertising.LEAdvertisement1Properties
}

func (a *Advertisement) DBusConn() *dbus.Conn {
	return a.conn
}

func (a *Advertisement) DBusObjectManager() *DBusObjectManager {
	return a.objectManager
}

func (a *Advertisement) DBusProperties() *DBusProperties {
	return a.iprops
}

func (a *Advertisement) GetProperties() bluez.Properties {
	return a.props
}

func (a *Advertisement) Path() dbus.ObjectPath {
	return a.path
}

func (a *Advertisement) Interface() string {
	return advertising.LEAdvertisement1Interface
}

func NewAdvertisement(adapterID string, props *advertising.LEAdvertisement1Properties) (*Advertisement, error) {

	adv := new(Advertisement)

	adv.props = props
	adv.path = nextAdvertismentPath(adapterID)

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	adv.conn = conn

	om, err := NewDBusObjectManager(conn)
	if err != nil {
		return nil, err
	}
	adv.objectManager = om

	iprops, err := NewDBusProperties(conn)
	if err != nil {
		return nil, err
	}
	adv.iprops = iprops

	return adv, nil
}

// Expose to bluez an advertisment instance via the adapter advertisement manager
func ExposeAdvertisement(adapterID string, props *advertising.LEAdvertisement1Properties, discoverableTimeout uint32) (func(), error) {

	log.Tracef("Retrieving adapter instance %s", adapterID)
	a, err := GetAdapter(adapterID)
	if err != nil {
		return nil, err
	}

	adv, err := NewAdvertisement(adapterID, props)
	if err != nil {
		return nil, err
	}

	err = ExposeDBusService(adv)
	if err != nil {
		return nil, err
	}

	log.Debug("Setup adapter")
	err = a.SetDiscoverable(true)
	if err != nil {
		return nil, err
	}

	err = a.SetDiscoverableTimeout(discoverableTimeout)
	if err != nil {
		return nil, err
	}
	err = a.SetPowered(true)
	if err != nil {
		return nil, err
	}

	log.Trace("Registering LEAdvertisement1 instance")
	advManager, err := advertising.NewLEAdvertisingManager1FromAdapterID(adapterID)
	if err != nil {
		return nil, err
	}

	err = advManager.RegisterAdvertisement(adv.Path(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	cancel := func() {
		decreaseAdvertismentCounter()
		err := advManager.UnregisterAdvertisement(adv.Path())
		if err != nil {
			log.Warn(err)
		}
		err = a.SetProperty("Discoverable", false)
		if err != nil {
			log.Warn(err)
		}
	}

	return cancel, nil
}
