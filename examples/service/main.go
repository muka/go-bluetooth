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

func main() {

	log.SetLevel(log.DebugLevel)

	var err error

	agent, err := createAgent()
	if err != nil {
		log.Errorf("createAgent: %s", err)
		os.Exit(1)
	}

	defer agent.Release()

	app, err := registerApplication(serviceAdapterID)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	defer app.StopAdvertising()

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
