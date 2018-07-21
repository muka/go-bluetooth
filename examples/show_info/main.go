package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	log "github.com/sirupsen/logrus"
)

//ShowInfoExample show informations for hardcoded MiBand2 on hci0
func main() {

	// Load adapter and device info
	adapterID := "hci0"
	deviceID := "ED:4B:79:DC:D4:D4" // MI Band 2

	LoadInfoExample(adapterID, deviceID)

	devices, err := api.GetDevices()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info(devices)

}

//LoadInfoExample show basic informations regarding a device
func LoadInfoExample(adapterID string, deviceID string) {

	adapter := profile.NewAdapter1(adapterID)

	log.Info("Adapter info\n---")
	log.Infof("Name: %s\n", adapter.Properties.Name)
	log.Infof("Modalias: %s\n", adapter.Properties.Modalias)
	log.Infof("Devices: %s\n", adapter.Properties.UUIDs)

	device := profile.NewDevice1(
		fmt.Sprintf(
			"/org/bluez/%s/dev_%s",
			adapterID,
			strings.Replace(deviceID, ":", "_", -1),
		),
	)

	log.Info("Device info\n---")
	log.Infof("Name: %s\n", device.Properties.Name)
	log.Infof("Modalias: %s\n", device.Properties.Modalias)

}
