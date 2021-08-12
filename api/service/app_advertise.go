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

	for _, svc := range app.GetServices() {
		adv.ServiceUUIDs = append(adv.ServiceUUIDs, svc.UUID)
	}

	cancel, err := api.ExposeAdvertisement(app.adapterID, adv, timeout)
	return cancel, err
}
