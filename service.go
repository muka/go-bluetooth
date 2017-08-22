package main

import (
	"github.com/muka/go-bluetooth/service"
	"github.com/prometheus/common/log"
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
