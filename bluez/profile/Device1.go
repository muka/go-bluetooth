package profile

import (
	"github.com/muka/go-bluetooth/src/gen/profile/device"
)

type Device1 struct {
	*device.Device1
}

// NewDevice1 create a new Device1 client
func NewDevice1(objectPath string) (*Device1, error) {

	d1, err := device.NewDevice1(objectPath)
	if err != nil {
		return nil, err
	}

	dev := new(Device1)
	dev.Device1 = d1

	return dev, nil
}
