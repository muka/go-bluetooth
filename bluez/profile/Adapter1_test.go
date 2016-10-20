package profile

import (
	"testing"
)

func TestNewAdapter1(t *testing.T) {
	t.Log("Create Adapter1")

	a := NewAdapter1("hci0")

	t.Log("Start Discovery")
	err := a.StartDiscovery()
	if err != nil {
		t.Log("Error on StartDiscovery")
		t.Fatal(err)
	}

	t.Log("Stop Discovery")
	err = a.StopDiscovery()
	if err != nil {
		t.Log("Error on StartDiscovery")
		t.Fatal(err)
	}

	t.Skipped()

}
