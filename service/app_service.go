package service

import (
	"strconv"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

//CreateService create a new GattService1 instance
func (app *Application) CreateService(props *gatt.GattService1Properties, advertisedOptional ...bool) (*GattService1, error) {
	app.config.serviceIndex++
	appPath := string(app.Path())
	if appPath == "/" {
		appPath = ""
	}

	advertise := false
	if len(advertisedOptional) > 0 {
		advertise = advertisedOptional[0]
	}

	path := appPath + "/service" + strconv.Itoa(app.config.serviceIndex)
	c := &GattService1Config{
		app:        app,
		objectPath: dbus.ObjectPath(path),
		ID:         app.config.serviceIndex,
		conn:       app.config.conn,
		advertised: advertise,
	}
	s, err := NewGattService1(c, props)

	return s, err
}

//AddService add service to expose
func (app *Application) AddService(service *GattService1) error {

	app.services[service.Path()] = service

	err := service.Expose()
	if err != nil {
		return err
	}

	err = app.exportTree()
	if err != nil {
		return err
	}

	err = app.GetObjectManager().AddObject(service.Path(), service.Properties())
	if err != nil {
		return err
	}

	return err
}

//RemoveService remove an exposed service
func (app *Application) RemoveService(service *GattService1) error {
	if _, ok := app.services[service.Path()]; ok {

		delete(app.services, service.Path())
		err := app.GetObjectManager().RemoveObject(service.Path())

		//TODO: remove chars + descritptors too
		if err != nil {
			return err
		}

		err = app.exportTree()
		if err != nil {
			return err
		}
	}
	return nil
}

//GetServices return the registered services
func (app *Application) GetServices() map[dbus.ObjectPath]*GattService1 {
	return app.services
}
