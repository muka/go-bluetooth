package profile

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/src/gen/profile/advertising"
)

// Possible values:
// "tx-power"
// "appearance"
// "local-name"
const (
	SupportedIncludesTxPower    = "tx-power"
	SupportedIncludesAppearance = "appearance"
	SupportedIncludesLocalName  = "local-name"
)

// NewLEAdvertisingManager1 create a new LEAdvertisingManager1 client
func NewLEAdvertisingManager1(hostID string) (*advertising.LEAdvertisingManager1, error) {
	return advertising.NewLEAdvertisingManager1(fmt.Sprintf("/org/bluez/%s", hostID))
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
