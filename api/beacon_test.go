package api

import (
	"testing"

	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBeaconParse(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	dev, err := device.NewDevice("hci1", "70:C9:4E:58:AA:7E")
	if err != nil {
		t.Fatal(err)
	}

	isBeacon, beacon, err := NewBeacon(dev)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isBeacon)

	log.Debugf("%++v", beacon)

}
