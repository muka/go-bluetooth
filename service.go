package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/service"
)

const (
	adapterID  = "hci0"
	objectName = "org.bluez"
	objectPath = "/org/bluez/example/service"
)

// const (
// 	adapterID  = "hci0"
// 	objectName = "go.bluetooth"
// 	objectPath = "/"
// )

func main() {

	log.SetLevel(log.DebugLevel)

	cfg := &service.ApplicationConfig{
		ObjectName: objectName,
		ObjectPath: objectPath,
	}
	app, err := service.NewApplication(cfg)
	if err != nil {
		log.Errorf("Failed to initialize app: %s", err.Error())
		return
	}

	err = app.Run()
	if err != nil {
		log.Errorf("Failed to run: %s", err.Error())
		return
	}

	serviceProps := &profile.GattService1Properties{
		Primary: true,
		UUID:    app.GenerateUUID(),
	}

	service1, err := app.CreateService(serviceProps)
	if err != nil {
		log.Errorf("Failed to create service: %s", err.Error())
		return
	}

	err = app.AddService(service1)
	if err != nil {
		log.Errorf("Failed to add service: %s", err.Error())
		return
	}

	charProps := &profile.GattCharacteristic1Properties{
		UUID: app.GenerateUUID(),
		Flags: []string{
			bluez.FlagCharacteristicRead,
			bluez.FlagCharacteristicWrite,
		},
	}
	char, err := service1.CreateCharacteristic(charProps)
	if err != nil {
		log.Errorf("Failed to create char: %s", err.Error())
		return
	}

	err = service1.AddCharacteristic(char)
	if err != nil {
		log.Errorf("Failed to add char: %s", err.Error())
		return
	}

	descProps := &profile.GattDescriptor1Properties{
		UUID: app.GenerateUUID(),
		Flags: []string{
			bluez.FlagDescriptorRead,
			bluez.FlagDescriptorWrite,
		},
	}
	desc, err := char.CreateDescriptor(descProps)
	if err != nil {
		log.Errorf("Failed to create char: %s", err.Error())
		return
	}

	err = char.AddDescriptor(desc)
	if err != nil {
		log.Errorf("Failed to add desc: %s", err.Error())
		return
	}

	log.Info("Application started, waiting for connections")

	//Register Application
	gattManager, err := api.GetGattManager(adapterID)
	if err != nil {
		log.Errorf("Failed to get GattManager1: %s", err.Error())
		return
	}

	err = gattManager.RegisterApplication(app.Path(), map[string]interface{}{})
	if err != nil {
		log.Errorf("Failed to register application: %s", err.Error())
		return
	}

	// createClient(objectName, objectPath)

	select {}
}

func createClient(name string, path string) {

	om := profile.NewObjectManager(name, path)
	objs, err := om.GetManagedObjects()

	if err != nil {
		log.Errorf("Error getting objects: %s", err.Error())
		return
	}

	log.Infof("Got objects: %v", objs)
}
