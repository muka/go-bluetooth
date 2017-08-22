package main

import (
	"log"

	"github.com/muka/go-bluetooth/linux"
)

var adapterID = "hci0"

func main() {

	log.Printf("Starting adapter %s", adapterID)

	hciconfig := linux.HCIConfig{}
	res, err := hciconfig.Up()
	if err != nil {
		log.Printf("Failed to start device %s: %s", adapterID, err.Error())
		return
	}
	log.Printf("Address %s, enabled %t", res.Address, res.Enabled)

	res, err = hciconfig.Down()
	if err != nil {
		log.Printf("Failed to stop device %s: %s", adapterID, err.Error())
		return
	}
	log.Printf("Address %s, enabled %t", res.Address, res.Enabled)

}
