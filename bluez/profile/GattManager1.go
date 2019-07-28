package profile

import (
	"fmt"

	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

// NewGattManager1 create a new GattManager1 client
func NewGattManager1(hostID string) (*gatt.GattManager1, error) {
	return gatt.NewGattManager1(fmt.Sprintf("/org/bluez/%s", hostID))
}
