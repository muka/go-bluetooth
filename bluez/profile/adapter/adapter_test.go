package adapter

import (
	"fmt"
	"log"
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
)

func TestGetAdapterIDFromPath(t *testing.T) {

	hci := "hci999"
	path := fmt.Sprintf("%s/%s", bluez.OrgBluezPath, hci)
	adapterID, err := ParseAdapterID(dbus.ObjectPath(path))
	if err != nil {
		t.Fatal(err)
	}
	if hci != adapterID {
		t.Fatal(fmt.Errorf("%s != %s", hci, adapterID))
	}

	path += "/dev_AA_BB_CC"
	adapterID, err = ParseAdapterID(dbus.ObjectPath(path))
	if err != nil {
		t.Fatal(err)
	}
	if hci != adapterID {
		t.Fatal(fmt.Errorf("%s != %s", hci, adapterID))
	}

}

func TestGetAdapterIDFromPathFail(t *testing.T) {
	hci := "foo1234"
	path := fmt.Sprintf("%s/%s", bluez.OrgBluezPath, hci)
	_, err := ParseAdapterID(dbus.ObjectPath(path))
	if err == nil {
		t.Fatal("Expected error parsing hci device name")
	}

	path = "foo/hci1"
	_, err = ParseAdapterID(dbus.ObjectPath(path))
	if err == nil {
		t.Fatal("Expected error parsing bluez base path")
	}
}

func TestDefaultAdapterSetGet(t *testing.T) {
	testAdapterID := "hci1"
	SetDefaultAdapterID(testAdapterID)
	adapterID := GetDefaultAdapterID()
	if adapterID != testAdapterID {
		log.Fatal("Failed to set default adapter")
	}
}

func TestAdapterExists(t *testing.T) {

	adapterID := GetDefaultAdapterID()

	exists, err := AdapterExists(adapterID)
	if err != nil {
		t.Fatal(err)
	}

	if exists == false {
		t.Errorf("Expected %s to exists", adapterID)
		t.Fatal()
	}

}

func TestGetAdapterNotExists(t *testing.T) {
	_, err := GetAdapter("foobar")
	if err == nil {
		t.Fatal("adapter should not exists")
	}
}

func TestGetAdapterFromDevicePath(t *testing.T) {

	a := getDefaultAdapter(t)

	list, err := a.GetDeviceList()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Log("Cannot test GetAdapterFromDevicePath, empty device list")
		return
	}

	_, err = GetAdapterFromDevicePath(list[0])
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetAdapterFromDevicePathFail(t *testing.T) {
	_, err := GetAdapterFromDevicePath(dbus.ObjectPath("foo/test/hci1/dev_AA_BB_CC_DD_EE_FF"))
	if err == nil {
		t.Fatal("Expected failure parsing device path")
	}
}

func TestGetAdapter(t *testing.T) {
	a := getDefaultAdapter(t)
	if a.Properties.Address == "" {
		log.Fatal("Properties should not be empty")
	}

}
