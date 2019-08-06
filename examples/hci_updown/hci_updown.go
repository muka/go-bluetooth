package hci_updown_example

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/muka/go-bluetooth/hw/linux/hci"
	log "github.com/sirupsen/logrus"
)

//HciUpDownExample hciconfig up / down
func Run(rawAdapterID string) error {

	log.Info("Turn down")

	adapterID, err := strconv.Atoi(strings.Replace(rawAdapterID, "hci", "", 1))
	if err != nil {
		return err
	}

	err = hci.Down(adapterID)
	if err != nil {
		return fmt.Errorf("Failed to stop device hci%d: %s", adapterID, err.Error())
	}

	log.Info("Turn on")
	err = hci.Up(adapterID)
	if err != nil {
		return fmt.Errorf("Failed to start device hci%d: %s", adapterID, err.Error())
	}

	log.Info("Done.")

	return nil
}
