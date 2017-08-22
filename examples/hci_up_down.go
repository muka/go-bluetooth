package examples

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/linux"
)

//HciUpDownExample hciconfig up / down
func HciUpDownExample(adapterID string) {

	log.Infof("Starting adapter %s", adapterID)

	hciconfig := linux.HCIConfig{}
	res, err := hciconfig.Up()
	if err != nil {
		log.Errorf("Failed to start device %s: %s", adapterID, err.Error())
		return
	}
	log.Infof("Address %s, enabled %t", res.Address, res.Enabled)

	res, err = hciconfig.Down()
	if err != nil {
		log.Errorf("Failed to stop device %s: %s", adapterID, err.Error())
		return
	}
	log.Infof("Address %s, enabled %t", res.Address, res.Enabled)

}
