package main

import (
	"github.com/muka/go-bluetooth/linux"
	log "github.com/sirupsen/logrus"
)

//HciUpDownExample hciconfig up / down
func main() {

	log.Info("Turn down")

	adapterID := 0

	err := linux.Down(adapterID)
	if err != nil {
		log.Errorf("Failed to stop device hci%d: %s", adapterID, err.Error())
		return
	}

	log.Info("Turn on")
	err = linux.Up(adapterID)
	if err != nil {
		log.Errorf("Failed to start device hci%d: %s", adapterID, err.Error())
		return
	}

	log.Info("Done.")

	// adapterID := "hci0"
	//
	// log.Infof("Starting adapter %s", adapterID)
	//
	// hciconfig := linux.HCIConfig{}
	// res, err := hciconfig.Up()
	// if err != nil {
	// 	log.Errorf("Failed to start device %s: %s", adapterID, err.Error())
	// 	return
	// }
	// log.Infof("Address %s, enabled %t", res.Address, res.Enabled)
	//
	// res, err = hciconfig.Down()
	// if err != nil {
	// 	log.Errorf("Failed to stop device %s: %s", adapterID, err.Error())
	// 	return
	// }
	// log.Infof("Address %s, enabled %t", res.Address, res.Enabled)

}
