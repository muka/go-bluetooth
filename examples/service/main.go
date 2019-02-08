package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	serviceAdapterID = "hci0"
	clientAdapterID  = "hci1"

	objectName = "org.bluez"
	objectPath = "/go_bluetooth/example/service"

	appName = "go-bluetooth"
)

func main() {

	log.SetLevel(log.DebugLevel)

	err := createClient(clientAdapterID, objectName, objectPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = registerApplication(serviceAdapterID)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	select {}
}
