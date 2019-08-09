package api

import (
	"strings"
	"testing"

	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	eddystone "github.com/suapapa/go_eddystone"
)

func testNewBeacon(t *testing.T, frame eddystone.Frame) Beacon {

	dev := &device.Device1{
		Properties: &device.Device1Properties{
			Name:  "test_eddystone",
			UUIDs: []string{"FEAA"},
			ServiceData: map[string]interface{}{
				"FEAA": []byte(frame),
			},
		},
	}

	isBeacon, beacon, err := NewBeacon(dev)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isBeacon)
	assert.True(t, beacon.IsEddystone())
	assert.Equal(t, string(beacon.Type), string(BeaconTypeEddystone))
	assert.IsType(t, BeaconEddystone{}, beacon.GetEddystone())

	return beacon
}

func TestParseEddystoneUID(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	uid := "EDD1-EBEA-C04E-5DEF-A017"
	instanceUid := "0BDB-8753-9B67"
	txpower := 120
	frame, err := eddystone.MakeUIDFrame(
		strings.Replace(uid, "-", "", -1),
		strings.Replace(instanceUid, "-", "", -1),
		txpower,
	)
	if err != nil {
		t.Fatal(err)
	}

	beacon := testNewBeacon(t, frame)

	assert.Equal(t, uid, beacon.GetEddystone().UID)
	assert.Equal(t, instanceUid, beacon.GetEddystone().InstanceUID)
	assert.Equal(t, txpower, beacon.GetEddystone().CalibratedTxPower)
}

func TestParseEddystoneTLM(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	var batt uint16 = 1000
	var temp float32 = 25
	var advCnt uint32 = 10
	var secCnt uint32 = 50
	frame, err := eddystone.MakeTLMFrame(batt, temp, advCnt, secCnt)
	if err != nil {
		t.Fatal(err)
	}

	beacon := testNewBeacon(t, frame)
	e := beacon.GetEddystone()

	assert.Equal(t, batt, e.TLMBatteryVoltage)
	assert.Equal(t, temp, e.TLMTemperature)
	assert.Equal(t, advCnt, e.TLMAdvertisingPDU)
	assert.Equal(t, secCnt, e.TLMLastRebootedTime)

	// log.Debugf("%+v", e)

}
