package main

import (
	"fmt"
	"github.com/muka/bluez-client/bluez"
)

func main() {

	hci0 := bluez.NewAdapter1("0")

	fmt.Println("Started discovery")
	hci0.StartDiscovery()

	fmt.Println("Stop discovery")
	hci0.StopDiscovery()

	fmt.Println("RM dev")
	hci0.RemoveDevice("device")

}
