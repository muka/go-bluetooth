package beacon

import (
	"testing"

	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestParseIBeacon(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	uuid := "010203040506070809101112131415"
	major := uint16(999)
	minor := uint16(111)
	measuredPower := uint16(80)

	b1, err := CreateIBeacon(uuid, major, minor, measuredPower)
	if err != nil {
		t.Fatal(err)
	}

	frames := b1.GetFrames()

	dev := &device.Device1{
		Properties: &device.Device1Properties{
			Name: "test_ibeacon",
			ManufacturerData: map[uint16]interface{}{
				appleBit: frames,
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
