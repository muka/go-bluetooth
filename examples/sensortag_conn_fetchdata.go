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
	"github.com/muka/go-bluetooth/devices"
	logging "github.com/op/go-logging"
	"github.com/muka/go-bluetooth/emitter"
	"sync"
	"time"
)
var client MQTT.Client
var log = logging.MustGetLogger("examples")
var adapterID = "hci0"

func main(){

	manager := api.NewManager()
	error := manager.RefreshState()
	if error != nil {
		panic(error)
	}
	
	SensorTag()
}
func SensorTag() {

	//.....................AdapterExists..................................
	
	boo,err := api.AdapterExists("hci0") 
		if err != nil {
		panic(err)
	}
	log.Debug("AdapterExists: ",boo)

	//...................start discovery on adapterId hci0..................

	err = api.StartDiscoveryOn("hci0")
	if err != nil {
		panic(err)
	}
	// wait a moment for the device to be spawn
	time.Sleep(time.Second)
	//.....................get the list of discovered devices....................
	
	devarr,err := api.GetDevices()
	if err != nil {
		panic(err)
	}
	//log.Debug("devarr",devarr[0])
	len := len( devarr )
	log.Debug("length: ",len)
	
	//.....................get device properties.........(name,status-connected,paired,uuids,Address)..........
	
	for i := 0; i< len; i++ { 
		prop1,err := devarr[i].GetProperties()
		
			if err != nil {
				log.Fatal(err)
			}
		log.Debug("DeviceProperties -ADDRESS: ",prop1.Address)
	
		ConnectAndFetchSensorDetailAndData(prop1.Address)
	}
	wg.Wait()
	
}

func ConnectAndFetchSensorDetailAndData(tagAddress string){

	dev, err := api.GetDeviceByAddress(tagAddress)
	if err != nil {
		panic(err)
	}
	log.Debug("device (dev): ",dev)
	
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
	log.Debug("sensor name: ",name)
		
	name1 := sensorTag.Humidity.GetName()
	log.Debug("sensor name: ",name1)
		
	mpu := sensorTag.Mpu.GetName()
	log.Debug("sensor name: ",mpu)
		
	barometric := sensorTag.Barometric.GetName()
	log.Debug("sensor name: ",barometric)
		
	luxometer:= sensorTag.Luxometer.GetName()
	log.Debug("sensor name: ",luxometer)

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
		
	err = sensorTag.Temperature.StartNotify(tagAddress)
	if err != nil {
	 	panic(err)
	}

	err = sensorTag.Humidity.StartNotify(tagAddress)
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
		switch x.SensorType {
			
			case "pressure":
				
				log.Debug("************************pressure***************************************" )
				log.Debug("SensorType: ",x.SensorType)
				log.Debug("BarometericPressureValue: ",x.BarometericPressureValue)
				log.Debug("BarometericPressureUnit: ",x.BarometericPressureUnit)
				log.Debug("BarometericTempValue: ",x.BarometericTempValue)
				log.Debug("BarometericTempUnit: ",x.BarometericTempUnit)
				log.Debug("SensorId : ",x.SensorId )
					
			case "temperature":	
				
				log.Debug("***********************temperature****************************************" )
				log.Debug("SensorType: ",x.SensorType)
				log.Debug("AmbientTempValue: ",x.AmbientTempValue)
				log.Debug("AmbientTempUnit : ",x.AmbientTempUnit )
				log.Debug("ObjectTempValue : ",x.ObjectTempValue )
				log.Debug("ObjectTempUnit : ",x.ObjectTempUnit )
				log.Debug("SensorId : ",x.SensorId )
					
			case "humidity":	

				log.Debug("**************************humidity*************************************" )
				log.Debug("SensorType: ",x.SensorType)
				log.Debug("HumidityValue: ",x.HumidityValue)
				log.Debug("HumidityUnit: ",x.HumidityUnit)
				log.Debug("HumidityTempValue: ",x.HumidityTempValue)
				log.Debug("HumidityTempUnit: ",x.HumidityTempUnit)
				log.Debug("SensorId : ",x.SensorId )

			case "mpu":	
				
				log.Debug("***************************mpu************************************" )
				log.Debug("SensorType: ",x.SensorType)
				log.Debug("mpuGyroscopeValue: ",x.MpuGyroscopeValue)
				log.Debug("mpuGyroscopeUnit: ",x.MpuGyroscopeUnit)
				log.Debug("mpuAccelerometerValue: ",x.MpuAccelerometerValue)
				log.Debug("mpuAccelerometerUnit: ",x.MpuAccelerometerUnit)
				log.Debug("mpuMagnetometerValue: ",x.MpuMagnetometerValue)
				log.Debug("mpuMagnetometerUnit: ",x.MpuMagnetometerUnit)
				log.Debug("SensorId : ",x.SensorId )
	
			case "luxometer":	
				
				log.Debug("**************************luxometer*************************************" )
				log.Debug("SensorType: ",x.SensorType)
				log.Debug("LuxometerValue : ",x.LuxometerValue )
				log.Debug("LuxometerUnit: ",x.LuxometerUnit)
				log.Debug("SensorId : ",x.SensorId )
						
		}
	}))
		
	wg.Add (3)
		
}

var wg sync.WaitGroup