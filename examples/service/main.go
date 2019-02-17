package main

import (
	"os"

	"github.com/muka/go-bluetooth/api"
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

	var err error

	app, err := registerApplication(serviceAdapterID)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	adapter, err := api.GetAdapter(serviceAdapterID)
	if err != nil {
		log.Errorf("GetAadapter: %s", err)
		os.Exit(1)
	}

	adapterProps, err := adapter.GetProperties()
	if err != nil {
		log.Errorf("adapter.GetProperties: %s", err)
		os.Exit(1)
	}

	hwaddr := adapterProps.Address

	var serviceID string
	for _, service := range app.GetServices() {
		serviceID = service.GetProperties().UUID
		break
	}

	// w := 4
	// log.Infof("Waiting %dsec to start client ...", w)
	// time.Sleep(time.Second * time.Duration(w))
	// log.Info("Ok")

	err = createClient(clientAdapterID, hwaddr, serviceID)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	select {}
}
