package beacon

import (
	"testing"

	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestParseIBeacon(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	dev := &device.Device1{
		Properties: &device.Device1Properties{
			Name: "test_ibeacon",
			ManufacturerData: map[uint16]interface{}{
				0x76: []uint8{
					// type
					0x2, 0x15,
					// uuid
					0xb9, 0x40, 0x7f, 0x30, 0xf5, 0xf8, 0x46, 0x6e, 0xaf, 0xf9, 0x25, 0x55, 0x6b, 0x57, 0xfe, 0x6d,
					// mayor
					0x8d, 0x80,
					// minor
					0xe8, 0x48,
					// power
					0xb4,
				},
			},
		},
	}

	beacon, err := NewBeacon(dev)
	if err != nil {
		t.Fatal(err)
	}

	isBeacon := beacon.Parse()

	assert.True(t, isBeacon)
	assert.True(t, beacon.IsIBeacon())
	assert.Equal(t, string(beacon.Type), string(BeaconTypeIBeacon))
	assert.IsType(t, BeaconIBeacon{}, beacon.GetIBeacon())
}
