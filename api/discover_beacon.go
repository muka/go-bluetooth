package api

import (
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

type BeaconType string

const (
	BeaconEddystone = "eddystone"
	BeaconIBeacon   = "ibeacon"
)

type Beacon struct {
	Type   BeaconType
	Device *device.Device1
}

// Load beacon inforamtion if available
func (b *Beacon) Load() bool {

	props := b.Device.Properties

	if len(props.AdvertisingData) == 0 && len(props.ManufacturerData) == 0 {
		return false
	}

	log.Debugf("AdvertisingData %++v", b.Device.Properties.AdvertisingData)
	log.Debugf("ManufacturerData %++v", b.Device.Properties.ManufacturerData)

	return false
}

func NewBeacon(dev *device.Device1) (bool, Beacon, error) {
	b := Beacon{
		Device: dev,
	}
	return b.Load(), b, nil
}
