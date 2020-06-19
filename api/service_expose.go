package api

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

type ExposedDBusService interface {
	Path() dbus.ObjectPath
	Interface() string
	GetProperties() bluez.Properties
	// App() *App
	DBusProperties() *DBusProperties
	DBusObjectManager() *DBusObjectManager
	DBusConn() *dbus.Conn
	// ExportTree() error
}

func RemoveDBusService(s ExposedDBusService) error {

	err := s.DBusObjectManager().RemoveObject(s.Path())
	if err != nil {
		return err
	}

	return nil
}

// Expose
// - Interface() (Service, Char or Descr)
// - Properties interface
func ExposeDBusService(s ExposedDBusService) (err error) {

	conn := s.DBusConn()

	if conn == nil {
		conn, err = dbus.SystemBus()
		if err != nil {
			return err
		}
	}

	log.Tracef("Expose %s (%s)", s.Path(), s.Interface())
	err = conn.Export(s, s.Path(), s.Interface())
	if err != nil {
		return err
	}

	err = s.DBusProperties().AddProperties(s.Interface(), s.GetProperties())
	if err != nil {
		return err
	}

	log.Tracef("Expose Properties interface (%s)", s.Path())
	s.DBusProperties().Expose(s.Path())

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//Properties
			prop.IntrospectData,
			// Exposed service introspectable
			{
				Name:       s.Interface(),
				Methods:    introspect.Methods(s),
				Properties: s.DBusProperties().Introspection(s.Interface()),
			},
		},
	}

	err = conn.Export(
		introspect.NewIntrospectable(node),
		s.Path(),
		"org.freedesktop.DBus.Introspectable",
	)
	if err != nil {
		return err
	}

	// v, _ := intrsp.Introspect()
	// fmt.Println("service_expose introspection", v)

	return nil
}
