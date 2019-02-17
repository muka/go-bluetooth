package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

type SupportedIncludesType string

// Possible values:
// "tx-power"
// "appearance"
// "local-name"
const (
	SupportedIncludesTypeTxPower    SupportedIncludesType = "tx-power"
	SupportedIncludesTypeAppearance SupportedIncludesType = "appearance"
	SupportedIncludesTypeLocalName  SupportedIncludesType = "local-name"
)

// NewLEAdvertisingManager1 create a new LEAdvertisingManager1 client
func NewLEAdvertisingManager1(hostID string) *LEAdvertisingManager1 {
	a := new(LEAdvertisingManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.LEAdvertisingManager1",
			Path:  "/org/bluez/" + hostID,
			Bus:   bluez.SystemBus,
		},
	)
	return a
}

// LEAdvertisingManager1Properties exposed properties for LEAdvertisingManager1
type LEAdvertisingManager1Properties struct {
	ActiveInstances    byte
	SupportedInstances byte
	SupportedIncludes  []SupportedIncludesType
}

// LEAdvertisingManager1 client
type LEAdvertisingManager1 struct {
	client *bluez.Client
}

// Close the connection
func (a *LEAdvertisingManager1) Close() {
	a.client.Disconnect()
}

//RegisterAdvertisement add a new advertisement service
func (a *LEAdvertisingManager1) RegisterAdvertisement(advertisement string, options map[string]interface{}) error {
	return a.client.Call("RegisterAdvertisement", 0, dbus.ObjectPath(advertisement), options).Store()
}

//UnregisterAdvertisement drop an advertisement service
func (a *LEAdvertisingManager1) UnregisterAdvertisement(advertisement string) error {
	return a.client.Call("UnregisterAdvertisement", 0, dbus.ObjectPath(advertisement)).Store()
}
