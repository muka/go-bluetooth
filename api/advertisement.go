package api

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/service"
	log "github.com/sirupsen/logrus"
)

const baseAdvertismentPath = "/org/bluez/apps/advertisement%d"

var advertisingCount int = -1

func nextAdvertismentPath() dbus.ObjectPath {
	advertisingCount++
	return dbus.ObjectPath(fmt.Sprintf(baseAdvertismentPath, advertisingCount))
}

func decreaseAdvertismentCounter() {
	advertisingCount--
	if advertisingCount < -1 {
		advertisingCount = -1
	}
}

// Expose to bluez an advertisment instance via the adapter advertisement manager
func ExposeAdvertisement(adapterID string, props *advertising.LEAdvertisement1Properties, discoverableTimeout uint32) (func(), error) {

	log.Tracef("Retrieving adapter instance %s", adapterID)
	a, err := GetAdapter(adapterID)
	if err != nil {
		return nil, err
	}

	log.Trace("Connecting to DBus")
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	log.Trace("Creating LEAdvertisement1 instance")
	config, err := service.NewLEAdvertisement1Config(nextAdvertismentPath(), conn)
	if err != nil {
		return nil, err
	}

	adv, err := service.NewLEAdvertisement1(config, props)
	if err != nil {
		return nil, err
	}

	log.Trace("Exposing LEAdvertisement1 instance")
	err = adv.Expose()
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
		advManager.UnregisterAdvertisement(adv.Path())
		a.SetProperty("Discoverable", false)
	}

	return cancel, nil
}
