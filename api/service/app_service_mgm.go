package service

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

func (app *App) GetServices() map[dbus.ObjectPath]*Service {
	return app.services
}

func (app *App) NewService(uuid string) (*Service, error) {

	s := new(Service)
	s.UUID = app.GenerateUUID(uuid)

	s.app = app
	s.chars = make(map[dbus.ObjectPath]*Char)
	s.path = dbus.ObjectPath(fmt.Sprintf("%s/service%s", app.Path(), strings.Replace(s.UUID, "-", "_", -1)[:8]))
	s.Properties = NewGattService1Properties(s.UUID)

	iprops, err := api.NewDBusProperties(s.App().DBusConn())
	if err != nil {
		return nil, err
	}
	s.iprops = iprops

	return s, nil
}

func (app *App) AddService(s *Service) error {

	app.services[s.Path()] = s

	err := s.Expose()
	if err != nil {
		return err
	}

	err = app.DBusObjectManager().AddObject(s.Path(), map[string]bluez.Properties{
		s.Interface(): s.GetProperties(),
	})
	if err != nil {
		return err
	}

	log.Tracef("Added GATT Service UUID=%s %s", s.UUID, s.Path())

	return nil
}

//RemoveService remove an exposed service
func (app *App) RemoveService(service *Service) error {
	if _, ok := app.services[service.Path()]; !ok {
		return nil
	}

	for _, char := range service.GetChars() {
		err := service.RemoveChar(char)
		if err != nil {
			return err
		}
	}

	err := api.RemoveDBusService(service)
	if err != nil {
		return err
	}

	delete(app.services, service.Path())

	return nil
}
