package device

import "github.com/godbus/dbus/v5"

type SetsItem struct {
	Object dbus.ObjectPath
	Dict   map[string]byte
}
