package sensortag_temperature_example

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/devices/sensortag"
)

// example of reading temperature from a TI sensortag
func Run(tagAddress, adapterID string) error {

	a, err := api.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	dev, err := a.GetDeviceByAddress(tagAddress)
	if err != nil {
		return err
	}

	if dev == nil {
		return fmt.Errorf("Device %s not found", tagAddress)
	}

	err = dev.Connect()
	if err != nil {
		return err
	}

	sensorTag, err := sensortag.NewSensorTag(dev)
	if err != nil {
		return err
	}

	// var readTemperature = func() {
	// 	temp, err := sensorTag.Temperature.Read()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	log.Printf("Temperature %v°", temp)
	// }

	var notifyTemperature = func(fn func(temperature float64)) {
		sensorTag.Temperature.StartNotify()
		select {}
	}

	// readTemperature()
	notifyTemperature(func(t float64) {
		log.Infof("Temperature update: %f", t)
	})

	return nil
}
