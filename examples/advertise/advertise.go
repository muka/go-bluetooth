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
		Duration:     100,
		LocalName:    "goadv",
	}

	props.AddServiceUUID("180D", "180F")
	props.AddData(0x26, []byte{0x01, 0x01, 0x00})
	props.AddManifacturerData(0xfff, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	props.AddServiceData("9999", []byte{0x00, 0x01, 0x02, 0x03, 0x04})

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

	log.Debug("Ready")

	select {}
}
