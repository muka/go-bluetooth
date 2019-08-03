package beacon_example

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/service"
	log "github.com/sirupsen/logrus"
)

const advertismentPath = "/org/bluez/example/advertisement0"

func Run(beaconType, adapterID string) error {

	log.Debugf("Retrieving adapter instance %s", adapterID)
	adapter, err := adapter.NewAdapter1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	props := advertising.LEAdvertisement1Properties{
		Duration:            100,
		Timeout:             1,
		DiscoverableTimeout: 0,
	}

	// Based on src/bluez/test/example-advertisement
	// props.AddServiceUUID("180D", "180F")
	// props.AddData(0x26, []byte{0x01, 0x01, 0x00})
	// props.AddManifacturerData(0xfff, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	// props.AddServiceData("9999", []byte{0x00, 0x01, 0x02, 0x03, 0x04})

	if beaconType == "ibeacon" {
		err = iBeacon(&props)
	} else {
		err = eddystoneBeacon(&props)
	}
	if err != nil {
		return err
	}

	log.Debug("Connecting to DBus")
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	log.Debug("Creating LEAdvertisement1 instance")
	config, err := service.NewLEAdvertisement1Config(advertismentPath, conn)
	if err != nil {
		return err
	}

	adv, err := service.NewLEAdvertisement1(config, &props)
	if err != nil {
		return err
	}

	log.Debug("Exposing LEAdvertisement1 instance")
	err = adv.Expose()
	if err != nil {
		return err
	}

	log.Debug("Setup adapter")
	err = adapter.SetProperty("Powered", true)
	if err != nil {
		return err
	}
	err = adapter.SetProperty("Discoverable", true)
	if err != nil {
		return err
	}
	err = adapter.SetProperty("DiscoverableTimeout", uint32(0))
	if err != nil {
		return err
	}

	log.Debug("Registering LEAdvertisement1 instance")
	advManager, err := advertising.NewLEAdvertisingManager1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	err = advManager.RegisterAdvertisement(adv.Path(), map[string]interface{}{})
	if err != nil {
		return err
	}

	defer func() {
		advManager.UnregisterAdvertisement(adv.Path())
		adapter.SetProperty("Discoverable", false)
	}()

	log.Debugf("%s ready", beaconType)

	select {}
}

// Credits
// https://scribles.net/creating-ibeacon-using-bluez-example-code-on-raspberry-pi/
func iBeacon(props *advertising.LEAdvertisement1Properties) error {

	// props.Type = advertising.AdvertisementTypeBroadcast
	props.Type = advertising.AdvertisementTypePeripheral
	props.LocalName = "go_ibeacon"

	company_id := uint16(0x004C)
	payload := []uint8{
		// beacon_type
		0x2, 0x15,
		// uuid
		0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,

		// major
		// 0x11,
		// 0x22,

		// minor
		// 0x33,
		// 0x44,

		// tx_power at 1m
		// 0x50,
	}

	props.AddManifacturerData(company_id, payload)
	return nil
}

// Based on
// https://github.com/google/eddystone/tree/master/eddystone-url
func eddystoneBeacon(props *advertising.LEAdvertisement1Properties) error {

	props.LocalName = "goeddystone"
	props.Type = advertising.AdvertisementTypePeripheral

	props.AddServiceUUID("FEAA")
	// props.AddServiceData("FEAA", f)
	props.AddServiceData("FEAA", []uint8{
		0x10, /* frame type Eddystone-URL */
		0x00, /* Tx power at 0m */
		0x00, /* URL Scheme Prefix http://www. */
		'b',
		'l',
		'u',
		'e',
		'z',
		0x01, /* .org/ */
	})

	return nil
}
