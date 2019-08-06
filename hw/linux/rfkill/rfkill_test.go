package rfkill

import "testing"

var testAdapterID = "hci0"

func TestGetAdapterStatus(t *testing.T) {
	_, err := GetAdapterStatus(testAdapterID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestToggleAdapter(t *testing.T) {
	err := ToggleAdapter(testAdapterID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOnAdapter(t *testing.T) {
	err := TurnOnAdapter(testAdapterID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOffAdapter(t *testing.T) {
	err := TurnOffAdapter(testAdapterID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOnBluetooth(t *testing.T) {
	err := TurnOnBluetooth()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOffBluetooth(t *testing.T) {
	err := TurnOffBluetooth()
	if err != nil {
		t.Fatal(err)
	}
}

func TestToggleBluetooth(t *testing.T) {
	err := ToggleBluetooth()
	if err != nil {
		t.Fatal(err)
	}
}
