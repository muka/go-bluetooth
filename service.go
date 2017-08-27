package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/service"
)

const (
	objectName = "go.bluetooth"
	objectPath = "/"
)

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

	props := &profile.GattService1Properties{
		Primary: true,
		UUID:    app.GenerateUUID(),
	}

	err = app.Run()
	if err != nil {
		log.Errorf("Failed to run: %s", err.Error())
		return
	}

	service1, err := app.CreateService(props)
	if err != nil {
		log.Errorf("Failed to create service: %s", err.Error())
		return
	}
	err = app.AddService(service1)
	if err != nil {
		log.Errorf("Failed to add service: %s", err.Error())
		return
	}

	log.Info("Application started, waiting for connections")

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
