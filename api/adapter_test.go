package api

import (
	"log"
	"testing"
)

func TestDefaultAdapterSetGet(t *testing.T) {
	testAdapterID := "hci1"
	SetDefaultAdapter(testAdapterID)
	adapterID := GetDefaultAdapter()
	if adapterID != testAdapterID {
		log.Fatal("Failed to set default adapter")
	}
}

func TestAdapterExists(t *testing.T) {

	adapterID := GetDefaultAdapter()

	exists, err := AdapterExists(adapterID)
	if err != nil {
		t.Fatal(err)
	}

	if exists == false {
		t.Errorf("Expected %s to exists", adapterID)
		t.Fatal()
	}

}

func TestGetAdapter(t *testing.T) {

	adapterID := GetDefaultAdapter()

	a, err := GetAdapter(adapterID)
	if err != nil {
		t.Fatal(err)
	}

	if a.Properties.Address == "" {
		log.Fatal("Properties should not be empty")
	}

}
