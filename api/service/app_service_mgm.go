package service

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

func (app *App) GetServices() map[dbus.ObjectPath]*Service {
	return app.services
}

func (app *App) NewService(uuid string) (*Service, error) {

	if len(uuid) == 4 {
		uuid = fmt.Sprintf(DefaultUUID, uuid)
		log.Warnf("Using default UUID as base: %s", uuid)
	}

	s := new(Service)
	s.chars = make(map[dbus.ObjectPath]*Char)
	s.path = dbus.ObjectPath(fmt.Sprintf("%s/service_%s", app.Path(), strings.Replace(uuid, "-", "_", -1)))
	s.props = NewGattService1Properties(uuid)

	iprops, err := NewDBusProperties()
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

	err = app.exportTree()
	if err != nil {
		return err
	}

	err = app.objectManager.AddObject(s.Path(), map[string]bluez.Properties{
		s.Interface(): s.Properties(),
	})
	if err != nil {
		return err
	}

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

	err := RemoveDBusService(service)
	if err != nil {
		return err
	}

	delete(app.services, service.Path())

	return nil
}
