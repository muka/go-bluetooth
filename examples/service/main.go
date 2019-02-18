package main

import (
	"os"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/linux/btmgmt"
	log "github.com/sirupsen/logrus"
)

const (
	serviceAdapterID = "hci0"
	clientAdapterID  = "hci1"

	objectName      = "org.bluez"
	objectPath      = "/go_bluetooth/example/service"
	agentObjectPath = "/go_bluetooth/example/agent"

	appName = "go-bluetooth"
)

func reset() {

	// turn off/on
	btmgmt1 := btmgmt.NewBtMgmt(serviceAdapterID)
	err := btmgmt1.Reset()
	if err != nil {
		log.Warnf("Reset %s: %s", serviceAdapterID, err)
		os.Exit(1)
	}

	btmgmt2 := btmgmt.NewBtMgmt(clientAdapterID)
	err = btmgmt2.Reset()
	if err != nil {
		log.Warnf("Reset %s: %s", clientAdapterID, err)
		os.Exit(1)
	}

	time.Sleep(time.Millisecond * 500)

}

func fail(where string, err error) {
	if err != nil {
		log.Errorf("%s: %s", where, err)
		os.Exit(1)
	}
}

func main() {

	log.SetLevel(log.DebugLevel)

	var err error

	agent, err := createAgent()
	fail("createAgent", err)

	defer agent.Release()

	app, err := registerApplication(serviceAdapterID)
	fail("registerApplication", err)

	defer app.StopAdvertising()

	adapter, err := api.GetAdapter(serviceAdapterID)
	fail("GetAdapter", err)

	adapterProps, err := adapter.GetProperties()
	fail("GetProperties", err)

	hwaddr := adapterProps.Address

	var serviceID string
	for _, service := range app.GetServices() {
		serviceID = service.GetProperties().UUID
		break
	}

	err = createClient(clientAdapterID, hwaddr, serviceID)
	fail("createClient", err)

	select {}
}
