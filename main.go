package main

import (
	"fmt"
	"github.com/muka/bluez-client/bluez"
)

func main() {

	adapter := bluez.NewAdapter1("hci0")

	fmt.Println("GetProperties:")
	props, err := adapter.GetProperties()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", props.Name)
	fmt.Printf("Modalias: %s\n", props.Modalias)

	//
	// fmt.Println("Started discovery")
	// adapter.StartDiscovery()
	//
	// fmt.Println("Stop discovery")
	// adapter.StopDiscovery()
	//
	// fmt.Println("RM dev")
	// adapter.RemoveDevice("device")

}
