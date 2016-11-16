package main

import (
	"github.com/muka/go-bluetooth/linux"
	"github.com/tj/go-debug"
)

var dbg = debug.Debug("bluez:main")

var adapterID = "hci0"

func main() {

	hciconfig := linux.HCIConfig{}
	res, err := hciconfig.Up(adapterID)
	if err != nil {
		panic(err)
	}
	dbg("Address %s, enabled %t", res.Address, res.Enabled)

	res, err = hciconfig.Down(adapterID)
	if err != nil {
		panic(err)
	}
	dbg("Address %s, enabled %t", res.Address, res.Enabled)

}
