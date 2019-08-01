package hci_updown_example

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/muka/go-bluetooth/linux"
	log "github.com/sirupsen/logrus"
)

//HciUpDownExample hciconfig up / down
func Run(rawAdapterID string) error {

	log.Info("Turn down")

	adapterID, err := strconv.Atoi(strings.Replace(rawAdapterID, "hci", "", 1))
	if err != nil {
		return err
	}

	err = linux.Down(adapterID)
	if err != nil {
		return fmt.Errorf("Failed to stop device hci%d: %s", adapterID, err.Error())
	}

	log.Info("Turn on")
	err = linux.Up(adapterID)
	if err != nil {
		return fmt.Errorf("Failed to start device hci%d: %s", adapterID, err.Error())
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

	return nil
}
