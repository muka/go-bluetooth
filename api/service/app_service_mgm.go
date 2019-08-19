package service

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

func (app *App) GetServices() map[dbus.ObjectPath]*Service {
	return app.services
}

func (app *App) NewService() (*Service, error) {

	s := new(Service)
	s.ID = len(app.services) + 1000

	if app.baseUUID == "" {
		rndUUID, err := RandomUUID()
		if err != nil {
			return nil, err
		}
		app.baseUUID = "%08x" + rndUUID[8:]
		log.Tracef("Base UUID: %s", rndUUID)
	}

	uuid := fmt.Sprintf(app.baseUUID, s.ID)

	s.app = app
	s.chars = make(map[dbus.ObjectPath]*Char)
	s.path = dbus.ObjectPath(fmt.Sprintf("%s/service_%s", app.Path(), strings.Replace(uuid, "-", "_", -1)[:8]))
	s.Properties = NewGattService1Properties(uuid)

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

	// log.Tracef("Added GATT Service ID=%d %s", s.ID, s.Path())

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
