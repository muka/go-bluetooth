package examples

import (
	"github.com/muka/bluez-client/api"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("examples")

func main() {

	// Load adapter and device info
	adapterID := "hci0"
	deviceID := "ED:4B:79:DC:D4:D4" // MI Band 2
	LoadInfo(adapterID, deviceID)

	devices, err := api.GetDevices()
	if err != nil {
		panic(err)
	}
	log.Println(devices)

}
