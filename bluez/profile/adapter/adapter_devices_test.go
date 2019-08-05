package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDefaultAdapter(t *testing.T) *Adapter1 {
	a, err := GetDefaultAdapter()
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestGetDeviceByAddress(t *testing.T) {
	a := getDefaultAdapter(t)
	list, err := a.GetDeviceByAddress("foobar")
	if err != nil {
		t.Fatal(err)
	}
	assert.Empty(t, list)
}

func TestGetDevices(t *testing.T) {
	a := getDefaultAdapter(t)
	_, err := a.GetDevices()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDeviceList(t *testing.T) {
	a := getDefaultAdapter(t)
	_, err := a.GetDeviceList()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlushDevices(t *testing.T) {
	a := getDefaultAdapter(t)
	err := a.FlushDevices()
	if err != nil {
		t.Fatal(err)
	}
	list, err := a.GetDevices()
	if err != nil {
		t.Fatal(err)
	}
	assert.Zero(t, len(list))
}
