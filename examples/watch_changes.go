package examples

import (
	"strings"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
)

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
func WatchChangesExample() {
	if !loadDevice() {
		discoverDevice()
	}
	select {}
}

func loadDevice() bool {
	defer api.Exit()
	return loadDevices()
}

func discoverDevice() {

	defer api.Exit()

	// logger.Debugf("Reset bluetooth device")
	// err := api.ToggleBluetooth()
	// if err != nil {
	// 	panic(err)
	// }

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
		emitter.On("adapter", emitter.NewCallback(func(ev emitter.Event) {
			info := ev.GetData().(api.AdapterEvent)

			if info.Status == api.DeviceAdded {

				logger.Debugf("Adapter %s added\n", info.Name)
				discoverDevices(info.Name)

			} else {
				logger.Debugf("Adapter %s removed\n", info.Name)
			}
		}))
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
	api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {

		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device

		if dev == nil {
			dbg("Device removed!")
			return
		}

		filterDevice(dev)
	}))

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
	props := dev.Properties
	if props.Address == dumpAddress {
		logger.Debugf("Found %s [addr:%s], list profiles", props.Name, props.Address)
		connectProfiles(dev)
		return true
	}

	return false
}

func listProfiles(dev *api.Device) {

	logger.Debug("Connected")

	// var LOCK = false

	dev.On("char", emitter.NewCallback(func(ev emitter.Event) {

		charEvent := ev.GetData().(api.GattCharacteristicEvent)
		charProps := charEvent.Properties

		substr := strings.ToUpper(charProps.UUID[4:8])
		dbg("Check for char %s", substr)
		serviceName := sensorTagUUIDs[substr]

		if serviceName != "" {

			// if LOCK {
			// 	return
			// }

			logger.Debugf("Found char %s (%s : %s)", serviceName, substr, charEvent.Path)

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
			// 		logger.Debug("Message %v", msg)
			//
			// 	}
			// }()
			//
			// err = gattChar.StartNotify()
			// if err != nil {
			// 	logger.Errorf("StartNotify error %s: %v", serviceName, err)
			// }

			// opts := make(map[string]dbus.Variant)
			//
			// logger.Debugf("Reading value for %s", serviceName)
			// raw, err := gattChar.ReadValue(opts)
			//
			// if err != nil {
			// 	logger.Errorf("Error reading %s: %v", serviceName, err)
			// } else {
			// 	logger.Debugf("Raw data %s: %v", serviceName, raw)
			// }

		}

	}))

}

func connectProfiles(dev *api.Device) {

	props := dev.Properties

	logger.Debugf("Connecting device %s", props.Name)
	err := dev.Connect()
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
