package examples

import (
	"fmt"
	"strings"

	"github.com/muka/go-bluetooth/bluez"
	"github.com/op/go-logging"
)

//LoadInfo show basic informations regardinf a device
func LoadInfo(adapterID string, deviceID string) {

	var log = logging.MustGetLogger(fmt.Sprintf("example:%s:%s", adapterID, deviceID))

	adapter := bluez.NewAdapter1(adapterID)

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

	log.Println("Device info\n---")
	log.Printf("Name: %s\n", device.Properties.Name)
	log.Printf("Modalias: %s\n", device.Properties.Modalias)

}
