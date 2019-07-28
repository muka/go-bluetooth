package profile

import (
	"github.com/muka/go-bluetooth/src/gen/profile/device"
)

// NewDevice1 create a new Device1 client
func NewDevice1(objectPath string) (*device.Device1, error) {
	return device.NewDevice1(objectPath)
}
