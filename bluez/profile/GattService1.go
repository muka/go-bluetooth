package profile

import (
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

// NewGattService1 create a new GattService1 client
func NewGattService1(objectPath string) (*gatt.GattService1, error) {
	return gatt.NewGattService1(objectPath)
}
