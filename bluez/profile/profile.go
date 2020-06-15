package profile

import "github.com/godbus/dbus/v5"

// BluezApi is the shared interface for the Bluez API implmentation
type BluezApi interface {
	Path() dbus.ObjectPath
	Interface() string
	Close()
}
