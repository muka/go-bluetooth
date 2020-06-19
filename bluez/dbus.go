package bluez

import (
	"errors"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

//Properties dbus serializable struct
// Use struct tags to control how the field is handled by Properties interface
// Example: field `dbus:writable,emit,myCallback`
// See Prop in github.com/godbus/dbus/v5/prop for configuration details
// Options:
// - writable: set the property as writable (Set will updated it). Omit for read-only
// - emit|invalidates: emit PropertyChanged, invalidates emit without disclosing the value. Omit for read-only
// - callback: a callable function in the struct compatible with the signature of Prop.Callback. Omit for no callback
type Properties interface {
	ToMap() (map[string]interface{}, error)
	Lock()
	Unlock()
}

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
	Path  dbus.ObjectPath
	Bus   BusType
}

// CloseConnections close all open connection to DBus
func CloseConnections() (err error) {
	for _, conn := range conns {
		if conn != nil {
			err = conn.Close()
			if err != nil {
				log.Warnf("Close: %s", err)
			}
		}
	}
	conns = make([]*dbus.Conn, 2)
	return err
}

//GetConnection get a DBus connection
func GetConnection(connType BusType) (*dbus.Conn, error) {
	switch connType {
	case SystemBus:
		if conns[SystemBus] == nil {
			// c.logger.Debug("Connecting to SystemBus")
			conn, err := dbus.SystemBus()
			if err != nil {
				return nil, err
			}
			conns[SystemBus] = conn
		}
		return conns[SystemBus], nil
	case SessionBus:
		if conns[SessionBus] == nil {
			// c.logger.Debug("Connecting to SessionBus")
			conn, err := dbus.SessionBus()
			if err != nil {
				return nil, err
			}
			conns[SessionBus] = conn
		}
		return conns[SessionBus], nil
	default:
		return nil, errors.New("Unmanged DBus type code")
	}
}
