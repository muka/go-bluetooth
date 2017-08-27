package service

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/prop"
)

//DbusErr a generic dbus Error
var DbusErr = dbus.NewError("org.freedesktop.Dbus.Error", nil)

//DbusErrIfaceNotFound interface not found
var DbusErrIfaceNotFound = prop.ErrIfaceNotFound

//DbusErrPropNotFound property not found
var DbusErrPropNotFound = prop.ErrPropNotFound
