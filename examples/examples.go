package examples

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {

	switch os.Args[1] {
	case "info":
		ShowInfoExample()
		break

	case "sensortag":
		SensorTagTemperatureExample()
		break

	case "watch":
		WatchChangesExample()
		break

	case "hci":
		HciUpDownExample("hci0")
		break

	default:
		log.Info(`
Sample code, may need configuration to work
Usage:
  - info
  - sensortag
  - watch
  - hci
`)
		break

	}

}
