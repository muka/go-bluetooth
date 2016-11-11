package main

import (
	"time"

	"github.com/muka/bluez-client/api"
	"github.com/muka/bluez-client/emitter"
	"github.com/op/go-logging"
	"github.com/tj/go-debug"
)

var logger = logging.MustGetLogger("main")
var dbg = debug.Debug("bluez:main")

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

func main() {

	defer api.Exit()

	logger.Debugf("Turning ON bluetooth")
	err := api.ToggleBluetooth()
	if err != nil {
		panic(err)
	}

	// wait a moment for the device to be spawn
	time.Sleep(time.Second)

	go waitAdapter()

	select {}
}

func waitAdapter() {
	if exists, err := api.AdapterExists(adapterID); !exists {
		if err != nil {
			panic(err)
		}
		logger.Debug("Waiting for adapter hci0")
		emitter.On("adapter", func(ev emitter.Event) {
			info := ev.GetData().(api.AdapterEvent)

			if info.Status == api.DeviceAdded {

				logger.Debugf("Adapter %s added\n", info.Name)
				discoverDevices(info.Name)

			} else {
				logger.Debugf("Adapter %s removed\n", info.Name)
			}
		})
	} else {
		discoverDevices(adapterID)
	}

	select {}
}

func discoverDevices(adapterID string) {

	if deviceIsCached() {
		return
	}

	err := api.StartDiscovery()
	if err != nil {
		panic(err)
	}

	logger.Debugf("Started discovery")
	api.On("discovery", func(ev api.Event) {

		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device

		filterDevice(dev)
	})

}

func deviceIsCached() bool {
	return loadDevices()
}

func loadDevices() bool {

	var err error
	devices, err := api.GetDevices()
	if err != nil {
		panic(err)
	}

	dbg("Loaded devices %d", len(devices))

	for _, dev := range devices {
		if filterDevice(&dev) {
			return true
		}
	}

	return false
}

func filterDevice(dev *api.Device) bool {
	props, err := dev.GetProperties()
	if err != nil {
		panic(err)
	}

	if props.Address == dumpAddress {
		logger.Debugf("Found %s [addr:%s], list profiles", props.Name, props.Address)
		connectProfiles(dev)
		return true
	}

	return false
}

func listProfiles(dev *api.Device) {

	logger.Debug("Connected")

	dev.On("service", func(ev api.Event) {
		serviceEvent := ev.GetData().(api.GattServiceEvent)
		serviceProps := serviceEvent.Properties

		substr := serviceProps.UUID[4:8]
		dbg("Check for %s", substr)
		serviceName := sensorTagUUIDs[substr]

		if serviceName != "" {
			logger.Debug("Found service %s (%s)", serviceName, substr)
		}

	})

	logger.Debug("Done.")
}

func connectProfiles(dev *api.Device) {

	props, err := dev.GetProperties()
	if err != nil {
		panic(err)
	}

	logger.Debugf("Connecting device %s", props.Name)
	err = dev.Connect()
	if err != nil {
		panic(err)
	}

	// logger.Debugf("Pairing device %s", props.Name)
	// err = dev.Pair()
	// if err != nil {
	// 	panic(err)
	// }

	time.Sleep(time.Second * 5)
	listProfiles(dev)

}
