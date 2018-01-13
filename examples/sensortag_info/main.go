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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/devices"
	"github.com/muka/go-bluetooth/emitter"
)

var adapterID = "hci0"

func main() {

	manager := api.NewManager()
	error := manager.RefreshState()
	if error != nil {
		panic(error)
	}

	ShowSensorTagInfo("hci0")
}

//ShowSensorTagInfo show info from a sensor tag
func ShowSensorTagInfo(adapterID string) {

	//.....................AdapterExists..................................

	boo, err := api.AdapterExists(adapterID)
	if err != nil {
		panic(err)
	}
	log.Debug("AdapterExists: ", boo)

	//...................start discovery on adapterId hci0..................

	err = api.StartDiscoveryOn(adapterID)
	if err != nil {
		panic(err)
	}
	// wait a moment for the device to be spawn
	time.Sleep(time.Second)
	//.....................get the list of discovered devices....................

	devarr, err := api.GetDevices()
	if err != nil {
		panic(err)
	}
	//log.Debug("devarr",devarr[0])
	len := len(devarr)
	log.Debug("length: ", len)

	//.....................get device properties.........(name,status-connected,paired,uuids,Address)..........

	for i := 0; i < len; i++ {
		prop1, err := devarr[i].GetProperties()

		if err != nil {
			log.Fatal(err)
		}
		log.Debug("DeviceProperties -ADDRESS: ", prop1.Address)

		ConnectAndFetchSensorDetailAndData(prop1.Address)
	}

}

func ConnectAndFetchSensorDetailAndData(tagAddress string) {

	dev, err := api.GetDeviceByAddress(tagAddress)
	if err != nil {
		panic(err)
	}
	log.Debug("device (dev): ", dev)

	if dev == nil {
		panic("Device not found")
	}

	if !dev.IsConnected() {
		log.Debug("not connected")

		err = dev.Connect()
		if err != nil {
			log.Fatal(err)
		}

	} else {

		log.Debug("already connected")

	}

	sensorTag, err := devices.NewSensorTag(dev)
	if err != nil {
		panic(err)
	}

	//.........getname returns sensorName.................

	name := sensorTag.Temperature.GetName()
	log.Debug("sensor name: ", name)

	name1 := sensorTag.Humidity.GetName()
	log.Debug("sensor name: ", name1)

	mpu := sensorTag.Mpu.GetName()
	log.Debug("sensor name: ", mpu)

	barometric := sensorTag.Barometric.GetName()
	log.Debug("sensor name: ", barometric)

	luxometer := sensorTag.Luxometer.GetName()
	log.Debug("sensor name: ", luxometer)

	//...........read sensorTag info.............................

	devInfo, err := sensorTag.DeviceInfo.Read()
	if err != nil {
		panic(err)
	}
	log.Debug("FirmwareVersion: ", devInfo.FirmwareVersion)
	log.Debug("HardwareVersion: ", devInfo.HardwareVersion)
	log.Debug("Manufacturer: ", devInfo.Manufacturer)
	log.Debug("Model: ", devInfo.Model)

	//........StartNotify......................

	err = sensorTag.Temperature.StartNotify()
	if err != nil {
		panic(err)
	}

	err = sensorTag.Humidity.StartNotify()
	if err != nil {
		panic(err)
	}

	err = sensorTag.Mpu.StartNotify(tagAddress)
	if err != nil {
		panic(err)
	}

	err = sensorTag.Barometric.StartNotify(tagAddress)
	if err != nil {
		panic(err)
	}

	err = sensorTag.Luxometer.StartNotify(tagAddress)
	if err != nil {
		panic(err)
	}

	//......receive emitted data of different sensors.............

	dev.On("data", emitter.NewCallback(func(ev emitter.Event) {
		x := ev.GetData().(api.DataEvent)
		log.Debugf("%++v", x)
	}))

}
