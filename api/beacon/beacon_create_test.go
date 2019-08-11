package beacon

import (
	"testing"

	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/stretchr/testify/assert"
)

func TestCreateIBeacon(t *testing.T) {

	uuid := "AAAABBBBCCCCDDDDAAAABBBBCCCCDDDD"
	maj := uint16(32000)
	min := uint16(11000)
	txPwr := uint16(0xB3)

	b, err := CreateIBeacon(uuid, maj, min, txPwr)
	if err != nil {
		t.Fatal(err)
	}

	b.Device = &device.Device1{
		Properties: &device.Device1Properties{},
	}

	b.Device.Properties.ServiceData = b.props.ServiceData

	isBeacon := b.Parse()
	assert.True(t, isBeacon)
	assert.Equal(t, BeaconTypeIBeacon, string(b.Type))

	assert.Equal(t, uuid, b.GetIBeacon().ProximityUUID)
	assert.Equal(t, maj, b.GetIBeacon().Major)
	assert.Equal(t, min, b.GetIBeacon().Minor)
	assert.Equal(t, txPwr, b.GetIBeacon().MeasuredPower)

}

func TestCreateEddystoneURL(t *testing.T) {

	url := "http://example.com"
	b, err := CreateEddystoneURL(url, 99)
	if err != nil {
		t.Fatal(err)
	}

	b.Device = &device.Device1{
		Properties: &device.Device1Properties{},
	}

	b.Device.Properties.ManufacturerData = b.props.ManufacturerData
	b.Device.Properties.UUIDs = b.props.ServiceUUIDs

	isBeacon := b.Parse()

	assert.True(t, isBeacon)
	assert.True(t, b.IsEddystone())
	assert.Equal(t, string(b.Type), string(BeaconTypeEddystone))
	assert.IsType(t, BeaconEddystone{}, b.GetEddystone())

	assert.Equal(t, url, b.GetEddystone().URL)
}

func TestCreateEddystoneTLM(t *testing.T) {

	batt := uint16(89)
	b, err := CreateEddystoneTLM(batt, 10.0, 10, 10)
	if err != nil {
		t.Fatal(err)
	}

	b.Device = &device.Device1{
		Properties: &device.Device1Properties{},
	}

	b.Device.Properties.ManufacturerData = b.props.ManufacturerData
	b.Device.Properties.UUIDs = b.props.ServiceUUIDs

	isBeacon := b.Parse()

	assert.True(t, isBeacon)
	assert.True(t, b.IsEddystone())
	assert.Equal(t, string(b.Type), string(BeaconTypeEddystone))
	assert.IsType(t, BeaconEddystone{}, b.GetEddystone())

	assert.Equal(t, batt, b.GetEddystone().TLMBatteryVoltage)
}

func TestCreateEddystoneUID(t *testing.T) {

	nsUID := "AAAAAAAAAABBBBBBBBBB"
	b, err := CreateEddystoneUID(nsUID, "123456123456", 99)
	if err != nil {
		t.Fatal(err)
	}

	b.Device = &device.Device1{
		Properties: &device.Device1Properties{},
	}

	b.Device.Properties.ManufacturerData = b.props.ManufacturerData
	b.Device.Properties.UUIDs = b.props.ServiceUUIDs

	isBeacon := b.Parse()

	assert.True(t, isBeacon)
	assert.True(t, b.IsEddystone())
	assert.Equal(t, string(b.Type), string(BeaconTypeEddystone))
	assert.IsType(t, BeaconEddystone{}, b.GetEddystone())

	assert.Equal(t, nsUID, b.GetEddystone().UID)
}
