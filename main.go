package main

import (
	"github.com/muka/go-bluetooth/linux"
	"github.com/tj/go-debug"
)

var adapterID = "hci0"

func main() {

	var dbg = debug.Debug("bluez:main")

	hciconfig := linux.HCIConfig{}
	res, err := hciconfig.Up()
	if err != nil {
		panic(err)
	}
	dbg("Address %s, enabled %t", res.Address, res.Enabled)

	res, err = hciconfig.Down()
	if err != nil {
		panic(err)
	}
	dbg("Address %s, enabled %t", res.Address, res.Enabled)

}
