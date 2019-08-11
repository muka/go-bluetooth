package beacon

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/service"
	log "github.com/sirupsen/logrus"
)

const baseAdvertismentPath = "/org/bluez/example/advertisement%d"

var advertisingCount int = -1

func getAdvertismentPath() dbus.ObjectPath {
	advertisingCount++
	return dbus.ObjectPath(fmt.Sprintf(baseAdvertismentPath, advertisingCount))
}

func decreaseAdvertismentCounter() {
	advertisingCount--
	if advertisingCount < -1 {
		advertisingCount = -1
	}
}

// Expose the beacon
func (b *Beacon) Expose(adapterID string, timeout uint16) (func(), error) {

	log.Debugf("Retrieving adapter instance %s", adapterID)
	a, err := api.GetAdapter(adapterID)
	if err != nil {
		return nil, err
	}

	props := b.props
	props.Type = advertising.AdvertisementTypeBroadcast

	if b.Name != "" {
		props.LocalName = b.Name
	}
	// Duration is set to 2sec by default
	// Not sure if duration can be mapped to interval.
	// Duration: 1,
	props.Timeout = timeout

	log.Debug("Connecting to DBus")
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	log.Debug("Creating LEAdvertisement1 instance")
	config, err := service.NewLEAdvertisement1Config(getAdvertismentPath(), conn)
	if err != nil {
		return nil, err
	}

	adv, err := service.NewLEAdvertisement1(config, props)
	if err != nil {
		return nil, err
	}

	log.Debug("Exposing LEAdvertisement1 instance")
	err = adv.Expose()
	if err != nil {
		return nil, err
	}

	log.Debug("Setup adapter")
	err = a.SetDiscoverable(true)
	if err != nil {
		return nil, err
	}

	err = a.SetDiscoverableTimeout(uint32(timeout))
	if err != nil {
		return nil, err
	}
	err = a.SetPowered(true)
	if err != nil {
		return nil, err
	}

	log.Debug("Registering LEAdvertisement1 instance")
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
