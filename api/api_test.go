package api

import (
	"testing"

	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/stretchr/testify/assert"
)

func TestGetAdapterID(t *testing.T) {

	defaultAdapterID := adapter.GetDefaultAdapterID()

	adapter.SetDefaultAdapterID("foo")
	adapterID := GetDefaultAdapterID()

	if adapterID != "foo" {
		t.Fatalf("Wrong adapter ID: %s", adapterID)
	}

	adapter.SetDefaultAdapterID(defaultAdapterID)
	adapterID = GetDefaultAdapterID()

	if adapterID != defaultAdapterID {
		t.Fatalf("Wrong adapter ID: %s", adapterID)
	}

}

func TestGetAdapter(t *testing.T) {

	a1, err := GetDefaultAdapter()
	if err != nil {
		t.Fatal(err)
	}

	id := GetDefaultAdapterID()
	a2, err := GetAdapter(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, a1.Properties.Address, a2.Properties.Address)

	err = Exit()
	if err != nil {
		t.Fatal(err)
	}

}
