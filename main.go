package main

import (
	"fmt"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/util"
	"strings"
)

func main() {

	log := util.NewLogger("main")

	adapterID := "hci0"
	deviceID := "ED:4B:79:DC:D4:D4" // MI Band 2

	adapter := bluez.NewAdapter1(adapterID)

	fmt.Println("GetProperties:")
	_, err := adapter.GetProperties()

	if err != nil {
		panic(err)
	}

	log.Printf("Name: %s\n", adapter.Properties.Name)
	log.Printf("Modalias: %s\n", adapter.Properties.Modalias)
	log.Printf("Devices: %s\n", adapter.Properties.UUIDs)

	// fmt.Println("Started discovery")
	// adapter.StartDiscovery()
	//
	// fmt.Println("Stop discovery")
	// adapter.StopDiscovery()
	//
	// fmt.Println("RM dev")
	// adapter.RemoveDevice("device")

	miband2 := bluez.NewDevice1(
		fmt.Sprintf(
			"/org/bluez/%s/dev_%s",
			adapterID,
			strings.Replace(deviceID, ":", "_", -1),
		),
	)

	_, err = miband2.GetProperties()

	if err != nil {
		panic(err)
	}

	log.Printf("Name: %s\n", miband2.Properties.Name)
	log.Printf("Modalias: %s\n", miband2.Properties.Modalias)

}
