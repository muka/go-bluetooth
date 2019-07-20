package main

import (
	"github.com/muka/go-bluetooth/devices/zoeetree"
	log "github.com/sirupsen/logrus"
)

var address = "F3:AC:E9:41:7B:AE"

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	zt := zoeetree.NewZoeeTree(address)
	err := zt.Connect(10)
	if err != nil {
		return err
	}

	log.Debugf("Found device %s", zt.Device.Properties.Name)

	select {}
}
