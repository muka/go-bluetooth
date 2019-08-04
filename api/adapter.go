package api

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

var defaultAdapter = "hci0"

func SetDefaultAdapter(a string) {
	defaultAdapter = a
}

func GetDefaultAdapter() string {
	return defaultAdapter
}

// AdapterExists checks if an adapter is available
func AdapterExists(adapterID string) (bool, error) {

	manager, err := GetManager()
	if err != nil {
		return false, err
	}

	objects := manager.GetObjects()

	path := dbus.ObjectPath(fmt.Sprintf("%s/%s", bluez.OrgBluezPath, adapterID))
	_, exists := objects.Load(path)

	return exists, nil
}

// GetAdapter return an adapter object instance
func GetAdapter(adapterID string) (*adapter.Adapter1, error) {

	if exists, err := AdapterExists(adapterID); !exists {
		if err != nil {
			return nil, fmt.Errorf("AdapterExists: %s", err)
		}
		return nil, fmt.Errorf("Adapter %s not found", adapterID)
	}

	return adapter.NewAdapter1FromAdapterID(adapterID)
}

func GetAdapterFromDevicePath(path dbus.ObjectPath) (*adapter.Adapter1, error) {

	d, err := device.NewDevice1(string(path))
	if err != nil {
		return nil, fmt.Errorf("Failed to load device %s", path)
	}

	a, err := adapter.NewAdapter1(string(d.Properties.Adapter))
	if err != nil {
		return nil, err
	}

	return a, nil
}

// ParseAdapterIDFromDevicePath return the adapterID from a device dbus object path
func ParseAdapterIDFromDevicePath(devicePath dbus.ObjectPath) (string, error) {

	parts := strings.Split(string(devicePath), "/")
	parts = parts[2:3]

	if len(parts) > 0 && parts[0][:3] == "hci" {
		adapterID := parts[0]
		return adapterID, nil
	}

	return "", fmt.Errorf("Failed to parse adapterID from %s", devicePath)
}
