package util

import (
	"testing"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	struct1 := struct {
		ManufacturerData map[uint16]interface{}
	}{}

	val1 := map[uint16][]byte{
		0x00: {0x01, 0x02, 0x03},
	}

	map1 := map[string]dbus.Variant{
		"ManufacturerData": dbus.MakeVariant(val1),
		"Foo":              dbus.MakeVariant(val1),
	}

	err := MapToStruct(&struct1, map1)
	if err != nil {
		t.Fatal(err)
	}

	val2, ok := struct1.ManufacturerData[0x00]
	assert.True(t, ok)
	assert.Equal(t, val2.([]byte), val1[0x00])

}

func TestNestedStructToMap(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	type fooStruct struct {
		Field1 string
		Field2 string
	}
	type barStruct struct {
		Field  uint32
		Nested *fooStruct
	}

	struct1 := struct {
		BarData barStruct
	}{BarData: barStruct{Nested: &fooStruct{}}}

	nestedVal := map[string]dbus.Variant{
		"Field1": dbus.MakeVariant("val3-1"),
		"Field2": dbus.MakeVariant("val3-2"),
	}
	val1 := map[string]dbus.Variant{
		"Field":  dbus.MakeVariant(uint32(9)),
		"Nested": dbus.MakeVariant(nestedVal),
	}

	map1 := map[string]dbus.Variant{
		"BarData": dbus.MakeVariant(val1),
	}

	err := MapToStruct(&struct1, map1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, struct1.BarData.Field, val1["Field"].Value().(uint32))
	assert.Equal(t, struct1.BarData.Nested.Field2, nestedVal["Field2"].Value().(string))
}
