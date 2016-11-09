package bluez

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/tj/go-debug"
)

var dbg = debug.Debug("bluez:dbus")

//BusType a type of DBus connection
type BusType int

const (
	// SessionBus uses the session bus
	SessionBus BusType = iota
	// SystemBus uses the system bus
	SystemBus
)

var conns = make([]*dbus.Conn, 2)

// Config pass configuration to a DBUS client
type Config struct {
	Name  string
	Iface string
	Path  string
	Bus   BusType
}

//GetConnection get a DBus connection
func GetConnection(connType BusType) (*dbus.Conn, error) {
	dbg("Get connection %d", connType)
	switch connType {
	case SystemBus:
		{
			if conns[SystemBus] == nil {
				// c.logger.Debug("Connecting to SystemBus")
				conn, err := dbus.SystemBus()
				if err != nil {
					return nil, err
				}
				conns[SystemBus] = conn
			}
			return conns[SystemBus], nil
		}
	case SessionBus:
		{
			if conns[SessionBus] == nil {
				// c.logger.Debug("Connecting to SessionBus")
				conn, err := dbus.SessionBus()
				if err != nil {
					return nil, err
				}
				conns[SessionBus] = conn
			}
			return conns[SessionBus], nil
		}
	default:
		{
			panic(errors.New("Unmanged DBus type code"))
		}
	}
}
