//this example starts discovery on adapter
//after discovery process GetDevices method
//returns list of discovered devices
//then with the help of mac address
//connectivity starts
//once sensors are connected it will
//fetch sensor name,manufacturer detail,
//firmware version, hardware version, model
//and sensor data...

package obex_push_example

import (
	"sync"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/obex"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func Run(targetAddress, filePath, adapterID string) error {

	a, err := api.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	dev, err := a.GetDeviceByAddress(targetAddress)
	if err != nil {
		return err
	}

	log.Debugf("device %s (%s)", dev.Properties.Name, dev.Properties.Address)

	if dev == nil {
		return err
	}

	props, err := dev.GetProperties()
	if err != nil {
		return err
	}
	if !props.Paired {
		log.Debug("not paired")

		err = dev.Pair()
		if err != nil {
			return err
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
			return err
		}
	}
	if tries >= maxRetry {
		log.Fatal("Max tries reached")
	}

	log.Debug("Session created: ", sessionPath)

	obexSession := obex.NewObexSession1(sessionPath)
	sessionProps, err := obexSession.GetProperties()
	if err != nil {
		return err
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
		return err
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
			return err
		}
		transferedPercent := (100 / float64(transProps.Size)) * float64(transProps.Transferred)

		log.Debug("Progress    : ", transferedPercent)
	}

	obexClient.RemoveSession(sessionPath)
	log.Debug(sessionPath)

	return nil
}
