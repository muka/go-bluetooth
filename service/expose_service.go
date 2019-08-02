package service

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
)

type ExposableService interface {
	Path() dbus.ObjectPath
	Interface() string
	Properties() bluez.Properties
	Conn() *dbus.Conn
}

func ExposeService(s ExposableService) error {

	conn := s.Conn()

	err := conn.Export(s, s.Path(), s.Interface())
	if err != nil {
		return err
	}

	propInterface, err := NewProperties(conn)
	if err != nil {
		return err
	}

	err = propInterface.AddProperties(s.Interface(), s.Properties())
	if err != nil {
		return err
	}

	propInterface.Expose(s.Path())

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//Properties
			prop.IntrospectData,
			//LEAdvertisement1
			{
				Name:       s.Interface(),
				Methods:    introspect.Methods(s),
				Properties: propInterface.Introspection(s.Interface()),
			},
		},
	}

	fmt.Printf("\n\n%++v\n\n", propInterface.Introspection(s.Interface()))

	err = conn.Export(
		introspect.NewIntrospectable(node),
		s.Path(),
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	return nil
}
