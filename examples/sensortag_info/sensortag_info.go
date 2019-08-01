//this example starts discovery on adapter
//after discovery process GetDevices method
//returns list of discovered devices
//then with the help of mac address
//connectivity starts
//once sensors are connected it will
//fetch sensor name,manufacturer detail,
//firmware version, hardware version, model
//and sensor data...

package sensortag_info_example

import (
	"fmt"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/devices/sensortag"
	"github.com/muka/go-bluetooth/emitter"
	log "github.com/sirupsen/logrus"
)

func Run(address string) error {

	manager, err := api.NewManager()
	if err != nil {
		return err
	}

	err = manager.RefreshState()
	if err != nil {
		return err
	}

	return ConnectAndFetchSensorDetailAndData(address)
}

//
// //ShowSensorTagInfo show info from a sensor tag
// func ShowSensorTagInfo(adapterID string) error {
//
// 	boo, err := api.AdapterExists(adapterID)
// 	if err != nil {
// 		return err
// 	}
// 	log.Debugf("AdapterExists: %t", boo)
//
// 	err = api.StartDiscoveryOn(adapterID)
// 	if err != nil {
// 		return err
// 	}
// 	// wait a moment for the device to be spawn
// 	time.Sleep(time.Second)
//
// 	devarr, err := api.GetDevices()
// 	if err != nil {
// 		return err
// 	}
// 	//log.Debug("devarr",devarr[0])
// 	len := len(devarr)
// 	log.Debugf("length: %d", len)
//
// 	for i := 0; i < len; i++ {
// 		prop1, err := devarr[i].GetProperties()
// 		if err != nil {
// 			log.Fatalf("Cannot load properties of %s: %s", devarr[i].Path, err.Error())
// 			continue
// 		}
// 		log.Debugf("DeviceProperties - ADDRESS: %s", prop1.Address)
//
// 		err = ConnectAndFetchSensorDetailAndData(prop1.Address)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }

// ConnectAndFetchSensorDetailAndData load an show sensor data
func ConnectAndFetchSensorDetailAndData(tagAddress string) error {

	dev, err := api.GetDeviceByAddress(tagAddress)
	if err != nil {
		return err
	}

	if dev == nil {
		return fmt.Errorf("device %s not found", tagAddress)
	}
	log.Debugf("device (dev): %v", dev)

	if !dev.IsConnected() {
		log.Debug("not connected")
		err = dev.Connect()
		if err != nil {
			return err
		}
	}

	sensorTag, err := sensortag.NewSensorTag(dev)
	if err != nil {
		return err
	}

	name := sensorTag.Temperature.GetName()
	log.Debugf("sensor name: %s", name)

	name1 := sensorTag.Humidity.GetName()
	log.Debugf("sensor name: %s", name1)

	mpu := sensorTag.Mpu.GetName()
	log.Debugf("sensor name: %s", mpu)

	barometric := sensorTag.Barometric.GetName()
	log.Debugf("sensor name: %s", barometric)

	luxometer := sensorTag.Luxometer.GetName()
	log.Debugf("sensor name: %s", luxometer)

	devInfo, err := sensorTag.DeviceInfo.Read()
	if err != nil {
		return err
	}

	log.Debug("FirmwareVersion: ", devInfo.FirmwareVersion)
	log.Debug("HardwareVersion: ", devInfo.HardwareVersion)
	log.Debug("Manufacturer: ", devInfo.Manufacturer)
	log.Debug("Model: ", devInfo.Model)

	err = sensorTag.Temperature.StartNotify()
	if err != nil {
		return err
	}

	err = sensorTag.Humidity.StartNotify()
	if err != nil {
		return err
	}

	err = sensorTag.Mpu.StartNotify(tagAddress)
	if err != nil {
		return err
	}

	err = sensorTag.Barometric.StartNotify(tagAddress)
	if err != nil {
		return err
	}

	err = sensorTag.Luxometer.StartNotify(tagAddress)
	if err != nil {
		return err
	}

	err = dev.On("data", emitter.NewCallback(func(ev emitter.Event) {
		x := ev.GetData().(sensortag.SensorTagDataEvent)
		log.Debugf("data received: %++v", x)
	}))

	return err
}
