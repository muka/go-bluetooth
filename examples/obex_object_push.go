//this example starts discovery on adapter
//after discovery process GetDevices method
//returns list of discovered devices
//then with the help of mac address
//connectivity starts
//once sensors are connected it will
//fetch sensor name,manufacturer detail,
//firmware version, hardware version, model
//and sensor data...

package main

import (
	"github.com/muka/go-bluetooth/api"
	"github.com/op/go-logging"
	"sync"
	"time"
	"github.com/muka/go-bluetooth/bluez/profile/obex"
	"os"
)

var log = logging.MustGetLogger("examples")

func main() {
	manager := api.NewManager()
	error := manager.RefreshState()
	if error != nil {
		panic(error)
	}

	SendFile(os.Args[1], os.Args[2])
}

func SendFile(targetAddress string, filePath string) {
	dev, err := api.GetDeviceByAddress(targetAddress)
	if err != nil {
		panic(err)
	}
	log.Debug("device (dev): ", dev)

	if dev == nil {
		panic("Device not found")
	}

	props, err := dev.GetProperties()
	if !props.Paired {
		log.Debug("not paired")

		err = dev.Pair()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Debug("already paired")
	}

	sessionArgs := map[string]interface{}{}
	sessionArgs["Target"] = "opp"

	obexClient := obex.NewObexClient1()

	tries := 1
	maxRetry := 20
	var sessionPath string
	for tries < maxRetry {
		log.Debug("Create Session...")
		sessionPath, err = obexClient.CreateSession(targetAddress, sessionArgs)
		if err == nil {
			break
		}

		tries++

		if err != nil {
			log.Error(err)
		}
	}
	if tries >= maxRetry {
		log.Fatal("Max tries reached")
	}

	log.Debug("Session created: ", sessionPath)

	obexSession := obex.NewObexSession1(sessionPath)
	sessionProps, err := obexSession.GetProperties()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Source		: ", sessionProps.Source)
	log.Debug("Destination	: ", sessionProps.Destination)
	log.Debug("Channel		: ", sessionProps.Channel)
	log.Debug("Target		: ", sessionProps.Target)
	log.Debug("Root			: ", sessionProps.Root)

	log.Debug("Init transmission on ", sessionPath)
	obexObjectPush := obex.NewObjectPush1(sessionPath)
	log.Debug("Send File: ", filePath)
	transPath, transProps, err := obexObjectPush.SendFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Transmission initiated: ", transPath)
	log.Debug("Status      : ", transProps.Status)
	log.Debug("Session     : ", transProps.Session)
	log.Debug("Name        : ", transProps.Name)
	log.Debug("Type        : ", transProps.Type)
	log.Debug("Time        : ", transProps.Time)
	log.Debug("Size        : ", transProps.Size)
	log.Debug("Transferred : ", transProps.Transferred)
	log.Debug("Filename    : ", transProps.Filename)

	for transProps.Transferred < transProps.Size {
		time.Sleep(1 * time.Second)

		obexTransfer := obex.NewObexTransfer1(transPath)
		transProps, err = obexTransfer.GetProperties()
		if err != nil {
			log.Fatal(err)
		}
		transferedPercent := (100 / float64(transProps.Size)) * float64(transProps.Transferred)

		log.Debug("Progress    : ", transferedPercent)
	}

	obexClient.RemoveSession(sessionPath)
	log.Debug(sessionPath)

}

var wg sync.WaitGroup
