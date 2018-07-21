package main

import (
	"os"
	"strings"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/linux"
	log "github.com/sirupsen/logrus"
)

var adapterID = "hci0"
var dumpAddress = "B0:B4:48:C9:4B:01"

var sensorTagUUIDs = map[string]string{
	"AA01": "TemperatureData",
	"AA02": "TemperatureConfig",
	"AA03": "TemperaturePeriod",
	"AA11": "AccelerometerData",
	"AA12": "AccelerometerConfig",
	"AA13": "AccelerometerPeriod",
	"AA21": "HumidityData",
	"AA22": "HumidityConfig",
	"AA23": "HumidityPeriod",
	"AA31": "MagnetometerData",
	"AA32": "MagnetometerConfig",
	"AA33": "MagnetometerPeriod",
	"AA41": "BarometerData",
	"AA42": "BarometerConfig",
	"AA44": "BarometerPeriod",
	"AA43": "BarometerCalibration",
	"AA51": "GyroscopeData",
	"AA52": "GyroscopeConfig",
	"AA53": "GyroscopePeriod",
	"AA61": "TestData",
	"AA62": "TestConfig",
	"CCC1": "ConnectionParams",
	"CCC2": "ConnectionReqConnParams",
	"CCC3": "ConnectionDisconnReq",
	"FFC1": "OADImageIdentify",
	"FFC2": "OADImageBlock",
}

//WatchChangesExample example events receival
func main() {
	ok, err := loadDevice()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if !ok {
		err := discoverDevice()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
	select {}
}

func loadDevice() (bool, error) {
	defer api.Exit()
	return loadDevices()
}

func discoverDevice() error {

	defer api.Exit()

	log.Debugf("Reset bluetooth device")
	err := linux.NewBtMgmt(adapterID).Reset()
	if err != nil {
		return err
	}

	// wait a moment for the device to be spawn
	time.Sleep(time.Second)

	go func() {
		for {
			err := waitAdapter()
			if err != nil {
				log.Error(err)
				log.Warn("Error on discovery, restarting in 10sec")
				time.Sleep(time.Second * 10)
			}
		}
	}()

	select {}
}

func waitAdapter() error {
	if exists, err := api.AdapterExists(adapterID); !exists {
		if err != nil {
			return err
		}
		log.Debug("Waiting for adapter hci0")
		err = emitter.On("adapter", emitter.NewCallback(func(ev emitter.Event) {
			info := ev.GetData().(api.AdapterEvent)

			if info.Status == api.DeviceAdded {

				log.Debugf("Adapter %s added\n", info.Name)
				discoverDevices(info.Name)

			} else {
				log.Debugf("Adapter %s removed\n", info.Name)
			}
		}))
		if err != nil {
			return err
		}
	} else {
		err = discoverDevices(adapterID)
		if err != nil {
			return err
		}
	}

	select {}
}

func discoverDevices(adapterID string) error {

	cached, err := deviceIsCached()
	if err != nil {
		return err
	}
	if cached {
		return nil
	}

	err = api.StartDiscovery()
	if err != nil {
		return err
	}

	log.Debugf("Started discovery")
	err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {

		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device

		if dev == nil {
			return
		}

		filterDevice(dev)
	}))

	return err
}

func deviceIsCached() (bool, error) {
	return loadDevices()
}

func loadDevices() (bool, error) {

	var err error
	devices, err := api.GetDevices()
	if err != nil {
		return false, err
	}

	for _, dev := range devices {
		if filterDevice(&dev) {
			return true, nil
		}
	}

	return false, nil
}

func filterDevice(dev *api.Device) bool {
	props := dev.Properties
	if props.Address == dumpAddress {
		log.Debugf("Found %s [addr:%s], list profiles", props.Name, props.Address)
		connectProfiles(dev)
		return true
	}

	return false
}

func listProfiles(dev *api.Device) error {

	log.Debug("Connected")

	// var LOCK = false

	err := dev.On("char", emitter.NewCallback(func(ev emitter.Event) {

		charEvent := ev.GetData().(api.GattCharacteristicEvent)
		charProps := charEvent.Properties

		substr := strings.ToUpper(charProps.UUID[4:8])
		serviceName := sensorTagUUIDs[substr]

		if serviceName != "" {

			log.Debugf("Found char %s (%s : %s)", serviceName, substr, charEvent.Path)

			// gattChar := profile.NewGattCharacteristic1(charEvent.DevicePath)

			// ch, err := gattChar.Register()
			// if err != nil {
			// 	panic(err)
			// }
			//
			// go func() {
			// 	for {
			//
			// 		if ch == nil {
			// 			return
			// 		}
			//
			// 		msg := <-ch
			//
			// 		if msg == nil {
			// 			return
			// 		}
			//
			// 		log.Debug("Message %v", msg)
			//
			// 	}
			// }()
			//
			// err = gattChar.StartNotify()
			// if err != nil {
			// 	log.Errorf("StartNotify error %s: %v", serviceName, err)
			// }

			// opts := make(map[string]dbus.Variant)
			//
			// log.Debugf("Reading value for %s", serviceName)
			// raw, err := gattChar.ReadValue(opts)
			//
			// if err != nil {
			// 	log.Errorf("Error reading %s: %v", serviceName, err)
			// } else {
			// 	log.Debugf("Raw data %s: %v", serviceName, raw)
			// }

		}

	}))

	return err
}

func connectProfiles(dev *api.Device) error {

	props := dev.Properties

	log.Debugf("Connecting device %s", props.Name)
	err := dev.Connect()
	if err != nil {
		return err
	}

	// log.Debugf("Pairing device %s", props.Name)
	// err = dev.Pair()
	// if err != nil {
	// 	return err
	// }

	time.Sleep(time.Second * 5)
	listProfiles(dev)

	return nil
}
