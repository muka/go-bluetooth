package media

import "github.com/godbus/dbus/v5"

// Item map to array{objects, properties}
type Item struct {
	Object   dbus.ObjectPath
	Property map[string]interface{}
}
