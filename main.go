package main

import (
	"bytes"
	"encoding/binary"

	"github.com/muka/bluez-client/api"
	"github.com/muka/bluez-client/emitter"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

var adapterID = "hci0"

var DumpAddress = "B0:B4:48:C9:4B:01"

func main() {

	defer api.Exit()

	logger.Debugf("Turning OFF bluetooth")
	err := api.TurnOffBluetooth()
	if err != nil {
		panic(err)
	}

	go waitAdapter()

	logger.Debugf("Turning ON bluetooth")
	err = api.TurnOnBluetooth()
	if err != nil {
		panic(err)
	}

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

	var err error
	devices, err := api.GetDevices()
	if err != nil {
		panic(err)
	}

	for _, dev := range devices {
		props, err := dev.GetProperties()
		if err != nil {
			panic(err)
		}

		if props.Address == DumpAddress {
			logger.Debugf("Found %s [addr:%s]", props.Name, props.Address)
			connectProfiles(&dev)
		}
	}

	// hci0, err := api.GetAdapter(adapterID)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, device := range devices {
	// 	logger.Debug("Dropping %s", device.Path)
	// 	go hci0.RemoveDevice(device.Path)
	// }
	//
	// err = api.StopDiscovery()
	// if err != nil {
	// 	logger.Debug(err)
	// }
	//
	// logger.Debug("Starting discovery on adapter %s\n", adapterID)
	// err = api.StartDiscovery()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// emitter.On("discovery", func(ev emitter.Event) {
	// 	info := ev.GetData().(api.DiscoveredDeviceEvent)
	// 	if info.Status == api.DeviceAdded {
	//
	// 		name := info.Device.GetProperties().Name
	//
	// 		if name != "MI Band 2" {
	// 			return
	// 		}
	//
	// 		logger.Debug("Found device %s, connecting profiles", name)
	// 		connectProfiles(info.Device)
	//
	// 		// 		logger.Debug("Found device %s, watching for property change", info.Device.GetProperties().Name)
	// 		//
	// 		// 		info.Device.On("change", func(ev emitter.Event) {
	// 		// 			changed := ev.GetData().(api.PropertyChanged)
	// 		// 			logger.Debug("%s: set %s = %s", info.Device.GetProperties().Name, changed.Field, changed.Value)
	// 		// 		})
	// 		//
	// 	} else {
	// 		logger.Debug("Removed device %s", info.Path)
	// 	}
	// })

}

func connectProfiles(dev *api.Device) {

	logger.Debug("Loading properties")

	err := dev.Connect()
	if err != nil {
		panic(err)
	}

	logger.Debug("Connected")

	props, err := dev.GetProperties()
	if err != nil {
		panic(err)
	}

	logger.Debugf("Got %d GattServices", len(props.GattServices))
	for _, path := range props.GattServices {

		logger.Debugf("Get Gatt service %s", path)
		service := dev.GetService(string(path))
		serviceProps, err := service.GetProperties()
		if err != nil {
			logger.Fatal(err)
			continue
		}

		logger.Debug("Got service %s\n", serviceProps.UUID)
		for _, charpath := range serviceProps.Characteristics {

			logger.Debug("Get Gatt char %s", charpath)
			char := dev.GetChar(string(charpath))
			charProps, err := char.GetProperties()

			if err != nil {
				logger.Error(err)
				continue
			}

			logger.Debugf("Got char %s\n", charProps.UUID)
			b, err := char.ReadValue()
			if err != nil {
				logger.Error(err)
				continue
			}

			if len(b) == 0 {
				logger.Debug("Empty bytearray")
				continue
			}

			logger.Debug("Char value is: %v\n", b)
			logger.Debug("string: %s\n", b)

			var n uint64
			buf := bytes.NewReader(b)
			err = binary.Read(buf, binary.LittleEndian, &n)
			if err != nil {
				logger.Debug("num: ")
				logger.Debug(n)
			}

			// uint64val, err := binary.ReadUvarint(buf)
			// if err != nil {
			// 	logger.Debug("uint64: ")
			// 	logger.Debug(uint64val)
			// }

			logger.Debug("---\n ")
		}

	}

	logger.Debug("Done.")
}
