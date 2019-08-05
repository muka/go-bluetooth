package util

import (
	"testing"

	"github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	struct1 := struct {
		ManufacturerData map[uint16]interface{}
	}{}

	val1 := map[uint16][]byte{
		0x00: []byte{0x01, 0x02, 0x03},
	}

	map1 := map[string]dbus.Variant{
		"ManufacturerData": dbus.MakeVariant(val1),
	}

	err := MapToStruct(&struct1, map1)
	if err != nil {
		t.Fatal(err)
	}

	val2, ok := struct1.ManufacturerData[0x00]
	assert.True(t, ok)
	assert.Equal(t, val2.([]byte), val1[0x00])

}
