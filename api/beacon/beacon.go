package beacon

import (
	"context"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

const appleBit = 0x004C

type BeaconType string

const (
	BeaconTypeEddystone = "eddystone"
	BeaconTypeIBeacon   = "ibeacon"
)

type Beacon struct {
	Name        string
	iBeacon     BeaconIBeacon
	eddystone   BeaconEddystone
	props       *advertising.LEAdvertisement1Properties
	Type        BeaconType
	Device      *device.Device1
	propchanged chan *bluez.PropertyChanged
}

func NewBeacon(dev *device.Device1) (Beacon, error) {
	b := Beacon{
		Name:   "gobluetooth",
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

// WatchDeviceChanges watch for properties changes
func (b *Beacon) WatchDeviceChanges(ctx context.Context) (chan bool, error) {
	var err error
	b.propchanged, err = b.Device.WatchProperties()
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)

	go func() {
		for {
			select {
			case changed := <-b.propchanged:

				if changed == nil {
					ctx.Done()
					return
				}

				if changed.Name == "ManufacturerData" || changed.Name == "ServiceData" {
					ch <- b.Parse()
				}

				break
			case <-ctx.Done():
				b.propchanged <- nil
				close(ch)
				break
			}
		}
	}()

	return ch, nil
}

func (b *Beacon) UnwatchDeviceChanges() error {
	err := b.Device.UnwatchProperties(b.propchanged)
	if err != nil {
		return err
	}

	return nil
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
	var data interface{}
	if b.IsIBeacon() {
		data = b.props.ManufacturerData[appleBit].([]byte)
	} else {
		data = b.props.ServiceData[eddystoneSrvcUid].([]byte)
	}
	if dataBytes, ok := b.getBytesFromData(data); ok {
		return dataBytes
	}
	return nil
}

// Load beacon information if available
func (b *Beacon) Parse() bool {
	if b.Device != nil {

		props := b.Device.Properties
		if b.parserEddystone(props.UUIDs, props.ServiceData) {
			return true
		}
		if b.parserIBeacon(props.ManufacturerData) {
			return true
		}

	}

	if b.props != nil {
		props := b.props
		if b.parserEddystone(props.ServiceUUIDs, props.ServiceData) {
			return true
		}
		if b.parserIBeacon(props.ManufacturerData) {
			return true
		}
	}

	return false
}

func (b *Beacon) parserIBeacon(manufacturerData map[uint16]interface{}) bool {
	if len(manufacturerData) == 0 {
		return false
	}
	if frames, ok := manufacturerData[appleBit]; ok {
		if frameBytes, ok := b.getBytesFromData(frames); ok {
			if len(frameBytes) < 22 {
				return false
			}

			b.Type = BeaconTypeIBeacon
			b.iBeacon = b.ParseIBeacon(frameBytes)
			return true
		}
	}
	return false
}

func (b *Beacon) parserEddystone(UUIDs []string, serviceData map[string]interface{}) bool {
	for _, uuid := range UUIDs {
		// 0000feaa-
		srcUUID := uuid
		if len(uuid) > 8 {
			uuid = uuid[4:8]
		}

		if strings.ToUpper(uuid) == eddystoneSrvcUid {
			if data, ok := serviceData[srcUUID]; ok {
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

func (b *Beacon) getBytesFromData(data interface{}) ([]byte, bool) {
	if variant, ok := data.(dbus.Variant); ok {
		if variantBytes, ok := variant.Value().([]byte); ok {
			return variantBytes, true
		}
	} else if dataBytes, ok := data.([]byte); ok {
		return dataBytes, true
	}

	return nil, false
}
