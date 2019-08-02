package advertise_example

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/service"
	log "github.com/sirupsen/logrus"
)

const advertismentPath = "/org/bluez/example/advertisement0"

func Run(adapterID string) error {

	log.Debugf("Retrieving adapter instance %s", adapterID)
	adapter, err := adapter.NewAdapter1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	props := advertising.LEAdvertisement1Properties{
		Type:         advertising.AdvertisementTypePeripheral,
		Discoverable: true,
		Duration:     10,
		LocalName:    "goadv",
	}

	// Based on src/bluez/test/example-advertisement
	// props.AddServiceUUID("180D", "180F")
	// props.AddData(0x26, []byte{0x01, 0x01, 0x00})
	// props.AddManifacturerData(0xfff, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	// props.AddServiceData("9999", []byte{0x00, 0x01, 0x02, 0x03, 0x04})

	// Credits
	// https://scribles.net/creating-ibeacon-using-bluez-example-code-on-raspberry-pi/

	company_id := uint16(0x004C)
	payload := []uint8{
		// beacon_type
		0x02, 0x15,
		// uuid
		0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
		// major
		// 0x2,
		// 0x22,
		// // minor
		// 0x0,
		// 0x44,
		// // tx_power (?)
		// 0xB3,
	}

	props.AddManifacturerData(company_id, payload)

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

	log.Debug("Power on adapter")
	err = adapter.SetProperty("Powered", true)
	if err != nil {
		return err
	}

	log.Debug("Registering LEAdvertisement1 instance")
	advManager, err := advertising.NewLEAdvertisingManager1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	err = advManager.RegisterAdvertisement(adv.Path(), map[string]dbus.Variant{})
	if err != nil {
		return err
	}
	defer advManager.UnregisterAdvertisement(adv.Path())

	log.Debug("IBeacon ready")

	select {}
}
