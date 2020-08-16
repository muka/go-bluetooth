package service

import (
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
)

func (app *App) GetAdvertisement() *advertising.LEAdvertisement1Properties {
	return app.advertisement
}

func (app *App) Advertise(timeout uint32) (func(), error) {

	adv := app.GetAdvertisement()

	adv.Timeout = uint16(timeout)
	adv.Duration = uint16(timeout)

	serviceUUIDs := []string{}
	for serviceUUID := range app.GetServices() {
		serviceUUIDs = append(serviceUUIDs, string(serviceUUID))
	}

	cancel, err := api.ExposeAdvertisement(app.adapterID, adv, timeout)
	return cancel, err
}
