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
	"github.com/muka/go-bluetooth/bluez/profile/battery"
	"github.com/muka/go-bluetooth/devices/sensortag"
	log "github.com/sirupsen/logrus"
)

func Run(address, adapterID string) error {

	a, err := api.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	dev, err := a.GetDeviceByAddress(address)
	if err != nil {
		return err
	}

	if dev == nil {
		return fmt.Errorf("device %s not found", address)
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

	batt, err := battery.NewBattery1(sensorTag.Device1.Path())
	if err != nil {
		log.Errorf("Cannot load battery profile: %s", err)
	} else {
		perc, err1 := batt.GetPercentage()
		if err1 != nil {
			log.Errorf("Cannot load battery percentage: %s", err)
		} else {
			log.Debugf("Battery: %d%%", perc)
		}
	}

	err = sensorTag.Temperature.StartNotify()
	if err != nil {
		return err
	}

	err = sensorTag.Humidity.StartNotify()
	if err != nil {
		return err
	}

	err = sensorTag.Mpu.StartNotify(address)
	if err != nil {
		return err
	}

	err = sensorTag.Barometric.StartNotify(address)
	if err != nil {
		return err
	}

	err = sensorTag.Luxometer.StartNotify(address)
	if err != nil {
		return err
	}

	go func() {
		// sensortag.SensorTagDataEvent
		for data := range sensorTag.Data() {
			log.Debugf("data received: %++v", data)
		}
	}()

	log.Debug("Waiting for data")

	select {}
}
