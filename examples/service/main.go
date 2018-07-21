package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	adapterID       = "hci0"
	clientAdapterID = "hci1"

	objectName = "org.bluez"
	objectPath = "/org/bluez/example/service"

	appName = "go-bluetooth"
)

func main() {

	log.SetLevel(log.DebugLevel)

	if len(os.Args) > 0 && os.Args[len(os.Args)-1] == "client" {
		createClient(clientAdapterID, objectName, objectPath)
	} else {
		err := registerApplication()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	select {}
}
