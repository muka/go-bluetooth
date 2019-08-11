package beacon

import (
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

const appleBit = 0x004C

type BeaconType string

const (
	BeaconTypeEddystone = "eddystone"
	BeaconTypeIBeacon   = "ibeacon"
)

type Beacon struct {
	iBeacon   BeaconIBeacon
	eddystone BeaconEddystone
	Type      BeaconType
	Device    *device.Device1
}

func NewBeacon(dev *device.Device1) (Beacon, error) {
	b := Beacon{
		Device: dev,
	}
	return b, nil
}

// IsEddystone return if the type of beacon is eddystone
func (b *Beacon) IsEddystone() bool {
	return b.Type == BeaconTypeEddystone
}

// IsIBeacon return if the type of beacon is ibeacon
func (b *Beacon) IsIBeacon() bool {
	return b.Type == BeaconTypeIBeacon
}

// GetEddystone return eddystone beacon information
func (b *Beacon) GetEddystone() BeaconEddystone {
	return b.eddystone
}

// GetIBeacon return if the type of beacon is ibeacon
func (b *Beacon) GetIBeacon() BeaconIBeacon {
	return b.iBeacon
}

// GetFrames return the bytes content
func (b *Beacon) GetFrames() []byte {
	if b.IsIBeacon() {
		return b.Device.Properties.ManufacturerData[appleBit].([]byte)
	}
	return b.Device.Properties.ServiceData[eddystoneSrvcUid].([]byte)
}

// Load beacon inforamtion if available
func (b *Beacon) Parse() bool {

	props := b.Device.Properties

	// log.Debugf("beacon props %++v", props)

	if len(props.ManufacturerData) > 0 {
		if frames, ok := props.ManufacturerData[appleBit]; ok {
			// log.Debug("Found iBeacon")
			// log.Debugf("iBeacon data: %d", frames)
			b.Type = BeaconTypeIBeacon
			b.iBeacon = b.ParseIBeacon(frames.([]byte))
		}
		return true
	}

	for _, uuid := range props.UUIDs {
		if uuid == eddystoneSrvcUid {
			if data, ok := props.ServiceData[eddystoneSrvcUid]; ok {
				// log.Debug("Found Eddystone")
				b.Type = BeaconTypeEddystone
				// log.Debugf("Eddystone data: %d", data)
				b.eddystone = b.ParseEddystone(data.([]byte))
				return true
			}
		}
	}

	return false
}
