package service

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/sirupsen/logrus"
)

//StartAdvertising advertise information for a service
func (app *Application) StartAdvertising(deviceInterface string) error {

	if app.advertisement != nil || app.adMgr != nil {
		logrus.Debugf("Already advertising on %s", deviceInterface)
		return nil
	}

	path := fmt.Sprintf("/org/bluez/advertisement/%s", deviceInterface)

	logrus.Debugf("Registering service on %s (%s)", deviceInterface, path)

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

	if len([]byte(strings.Join(serviceUUIDs, ","))) > 31 {
		logrus.Warn("Advertisment limit of 31 bytes may have been exceeded. Consider not exposing all services IDs with `CreateService(props, false)`")
	}

	props := &advertising.LEAdvertisement1Properties{
		Type:         advertising.AdvertisementTypePeripheral,
		LocalName:    app.config.LocalName,
		ServiceUUIDs: serviceUUIDs,
		// Appearance:   0,
	}

	var err error

	advertisement, err := NewLEAdvertisement1(config, props)
	if err != nil {
		return fmt.Errorf("NewLEAdvertisement1: %s", err)
	}

	app.advertisement = advertisement

	err = advertisement.Expose()
	if err != nil {
		app.advertisement = nil
		return fmt.Errorf("Expose: %s", err)
	}

	options := make(map[string]interface{})

	adMgr, err := advertising.NewLEAdvertisingManager1FromAdapterID(deviceInterface)
	if err != nil {
		return err
	}

	app.adMgr = adMgr

	err = app.adMgr.RegisterAdvertisement(dbus.ObjectPath(path), options)
	if err != nil {
		app.advertisement = nil
		app.adMgr = nil
		return fmt.Errorf("RegisterAdvertisement: %s", err)
	}

	adapter, err := adapter.NewAdapter1FromAdapterID(deviceInterface)
	if err != nil {
		return err
	}

	err = adapter.SetProperty("Discoverable", true)
	if err != nil {
		return err
	}

	err = adapter.SetProperty("Powered", true)
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

	err := app.advertisement.Release()
	if err != nil {
		return err
	}

	// err = app.adMgr.UnregisterAdvertisement(app.advertisement.config.objectPath)
	// if err != nil {
	// 	return err
	// }

	app.advertisement = nil
	app.adMgr = nil

	return err
}
