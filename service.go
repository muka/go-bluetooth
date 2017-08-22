package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/service"
)

func main() {

	cfg := &service.ApplicationConfig{}
	app, err := service.NewApplication(cfg)
	if err != nil {
		log.Errorf("Failed to initialize app: %s", err.Error())
		return
	}

	// err := app.AddService()

	err = app.Run()
	if err != nil {
		log.Errorf("Failed to run: %s", err.Error())
		return
	}

	log.Error("Application started, waiting for connections")
	select {}
}
