package main

import (
	"fmt"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/service"
	log "github.com/sirupsen/logrus"
)

func registerApplication(adapterID string) (*service.Application, error) {

	cfg := &service.ApplicationConfig{
		UUIDSuffix: "-0000-1000-8000-00805F9B34FB",
		UUID:       "1234",
		ObjectName: objectName,
		ObjectPath: objectPath,
		LocalName:  "GoBleSrvc",
	}
	app, err := service.NewApplication(cfg)
	if err != nil {
		log.Errorf("Failed to initialize app: %s", err.Error())
		return nil, err
	}

	err = app.Run()
	if err != nil {
		log.Errorf("Failed to run: %s", err.Error())
		return nil, err
	}

	err = exposeService(
		app,
		app.GenerateUUID("1111"),
		app.GenerateUUID("1111"),
		app.GenerateUUID("1111"),
		true,
	)
	if err != nil {
		return nil, err
	}

	err = exposeService(
		app,
		app.GenerateUUID("2222"),
		app.GenerateUUID("2222"),
		app.GenerateUUID("2222"),
		false,
	)
	if err != nil {
		return nil, err
	}

	log.Info("Application started, waiting for connections")

	//Register Application
	gattManager, err := api.GetGattManager(adapterID)
	if err != nil {
		return nil, fmt.Errorf("GetGattManager: %s", err)
	}

	err = gattManager.RegisterApplication(app.Path(), map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("RegisterApplication: %s", err.Error())
	}

	// Register our advertisement
	err = app.StartAdvertising(adapterID)
	if err != nil {
		return nil, fmt.Errorf("StartAdvertising: %s", err)
	}

	log.Info("Application registered and advertising.")
	return app, nil
}

func exposeService(
	app *service.Application,
	serviceUUID, characteristicUUID, descriptorUUID string,
	advertise bool,
) error {

	serviceProps := &profile.GattService1Properties{
		Primary: true,
		UUID:    serviceUUID,
	}

	// Set this service to be advertised
	service1, err := app.CreateService(serviceProps, advertise)
	if err != nil {
		log.Errorf("Failed to create service: %s", err.Error())
		return err
	}

	err = app.AddService(service1)
	if err != nil {
		log.Errorf("Failed to add service: %s", err.Error())
		return err
	}

	charProps := &profile.GattCharacteristic1Properties{
		UUID: characteristicUUID,
		Flags: []string{
			bluez.FlagCharacteristicRead,
			bluez.FlagCharacteristicWrite,
		},
	}
	char, err := service1.CreateCharacteristic(charProps)
	if err != nil {
		log.Errorf("Failed to create char: %s", err.Error())
		return err
	}

	err = service1.AddCharacteristic(char)
	if err != nil {
		log.Errorf("Failed to add char: %s", err.Error())
		return err
	}

	descProps := &profile.GattDescriptor1Properties{
		UUID: descriptorUUID,
		Flags: []string{
			bluez.FlagDescriptorRead,
			bluez.FlagDescriptorWrite,
		},
	}
	desc, err := char.CreateDescriptor(descProps)
	if err != nil {
		log.Errorf("Failed to create char: %s", err.Error())
		return err
	}

	err = char.AddDescriptor(desc)
	if err != nil {
		log.Errorf("Failed to add desc: %s", err.Error())
		return err
	}

	return nil
}
