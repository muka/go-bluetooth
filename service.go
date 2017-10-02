package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/service"
)

const (
	adapterID       = "hci0"
	clientAdapterID = "hci1"
	objectName      = "org.bluez"
	objectPath      = "/org/bluez/example/service"
)

func main() {

	log.SetLevel(log.DebugLevel)

	registerApplication()

	// createClient(objectName, objectPath)

	select {}
}

func createClient(name string, path string) {

	log.Info("Discovering devices")

	adapter := profile.NewAdapter1(clientAdapterID)

	err := adapter.StartDiscovery()
	if err != nil {
		log.Errorf("Failed to start discovery: %s", err.Error())
	}

	api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {

		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device

		if dev == nil {
			log.Infof("Device removed %s", dev.Path)
			return
		}

		log.Infof("Device found %s", dev.Path)

	}))

}

func registerApplication() {

	cfg := &service.ApplicationConfig{
		UUIDSuffix: "-0000-1000-8000-00805F9B34FB",
		UUID:       "1234",
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
		UUID:    app.GenerateUUID("2233"),
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
		UUID: app.GenerateUUID("3344"),
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
		UUID: app.GenerateUUID("4455"),
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

	adapter := profile.NewAdapter1(adapterID)
	err = adapter.SetProperty("Discoverable", dbus.MakeVariant(true))
	if err != nil {
		log.Errorf("Failed to set adapter %s discoverable: %s", adapterID, err.Error())
		return
	}

	log.Info("Application registered.")

}
