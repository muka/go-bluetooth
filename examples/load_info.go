package examples

import (
	"fmt"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/util"
	"strings"
)

//LoadInfo show basic informations regardinf a device
func LoadInfo(adapterID string, deviceID string) {

	var log = util.NewLogger(fmt.Sprintf("example:%s:%s", adapterID, deviceID))

	adapter := bluez.NewAdapter1(adapterID)

	_, err := adapter.GetProperties()
	if err != nil {
		panic(err)
	}

	log.Println("Adapter info\n---")
	log.Printf("Name: %s\n", adapter.Properties.Name)
	log.Printf("Modalias: %s\n", adapter.Properties.Modalias)
	log.Printf("Devices: %s\n", adapter.Properties.UUIDs)

	device := bluez.NewDevice1(
		fmt.Sprintf(
			"/org/bluez/%s/dev_%s",
			adapterID,
			strings.Replace(deviceID, ":", "_", -1),
		),
	)

	_, err = device.GetProperties()

	if err != nil {
		panic(err)
	}

	log.Println("Device info\n---")
	log.Printf("Name: %s\n", device.Properties.Name)
	log.Printf("Modalias: %s\n", device.Properties.Modalias)

}
