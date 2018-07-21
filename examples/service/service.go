package main

import (
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/service"
	log "github.com/sirupsen/logrus"
)

func registerApplication() error {

	cfg := &service.ApplicationConfig{
		UUIDSuffix: "-0000-1000-8000-00805F9B34FB",
		UUID:       "1234",
		ObjectName: objectName,
		ObjectPath: objectPath,
		LocalName:  "GolangService",
	}
	app, err := service.NewApplication(cfg)
	if err != nil {
		log.Errorf("Failed to initialize app: %s", err.Error())
		return err
	}

	err = app.Run()
	if err != nil {
		log.Errorf("Failed to run: %s", err.Error())
		return err
	}

	serviceProps := &profile.GattService1Properties{
		Primary: true,
		UUID:    app.GenerateUUID("2233"),
	}

	// Set this service to be advertised
	service1, err := app.CreateService(serviceProps, true)
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
		UUID: app.GenerateUUID("3344"),
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
		UUID: app.GenerateUUID("4455"),
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

	log.Info("Application started, waiting for connections")

	//Register Application
	gattManager, err := api.GetGattManager(adapterID)
	if err != nil {
		log.Errorf("Failed to get GattManager1: %s", err.Error())
		return err
	}

	err = gattManager.RegisterApplication(app.Path(), map[string]interface{}{})
	if err != nil {
		log.Errorf("Failed to register application: %s", err.Error())
		return err
	}

	// Register our advertisement
	err = app.StartAdvertising(adapterID)

	log.Info("Application registered and advertising.")
	return nil
}
