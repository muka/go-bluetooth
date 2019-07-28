package profile

import (
	"fmt"

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
