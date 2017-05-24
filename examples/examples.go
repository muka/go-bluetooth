package examples

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("examples")

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

	default:
		log.Info("Sample code, may need configuration to work\n\nUsage: \n - info\n - sensortag\n - watch\n")
		break

	}

}
