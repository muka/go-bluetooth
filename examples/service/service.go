package main

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/service"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
	log "github.com/sirupsen/logrus"
)

func registerApplication(adapterID string) (*service.Application, error) {

	cfg := &service.ApplicationConfig{
		UUIDSuffix: "-0000-1000-8000-00805F9B34FB",
		UUID:       "AAAA",
		ObjectName: objectName,
		ObjectPath: objectPath,
		LocalName:  "gobluetooth",
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

	err = gattManager.RegisterApplication(app.Path(), map[string]dbus.Variant{})
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

	serviceProps := service.NewGattService1Properties(serviceUUID)

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

	charProps := &service.GattCharacteristic1Properties{
		GattCharacteristic1Properties: &gatt.GattCharacteristic1Properties{
			UUID: characteristicUUID,
			Flags: []string{
				profile.FlagCharacteristicRead,
				profile.FlagCharacteristicWrite,
			},
		},
	}
	log.Debugf("%++v", charProps)
	char, err := service1.CreateCharacteristic(charProps)
	if err != nil {
		log.Errorf("Failed to create char: %s", err)
		return err
	}

	err = service1.AddCharacteristic(char)
	if err != nil {
		log.Errorf("Failed to add char: %s", err)
		return err
	}

	descProps := &service.GattDescriptor1Properties{
		GattDescriptor1Properties: &gatt.GattDescriptor1Properties{
			UUID: descriptorUUID,
			Flags: []string{
				profile.FlagDescriptorRead,
				profile.FlagDescriptorWrite,
			},
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
