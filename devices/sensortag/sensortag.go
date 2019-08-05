package sensortag

import (
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	log "github.com/sirupsen/logrus"
)

// DefaultRetry times
const DefaultRetry = 3

// DefaultRetryWait in millis
const DefaultRetryWait = 500

var dataChannel chan dbus.Signal

var sensorTagUUIDs = map[string]string{

	"TemperatureData":   "AA01",
	"TemperatureConfig": "AA02",
	"TemperaturePeriod": "AA03",

	"AccelerometerData":   "AA11",
	"AccelerometerConfig": "AA12",
	"AccelerometerPeriod": "AA13",

	"HumidityData":   "AA21",
	"HumidityConfig": "AA22",
	"HumidityPeriod": "AA23",

	"MagnetometerData":   "AA31",
	"MagnetometerConfig": "AA32",
	"MagnetometerPeriod": "AA33",

	"BarometerData":        "AA41",
	"BarometerConfig":      "AA42",
	"BarometerPeriod":      "AA44",
	"BarometerCalibration": "AA43",

	"GyroscopeData":   "AA51",
	"GyroscopeConfig": "AA52",
	"GyroscopePeriod": "AA53",

	"TestData":   "AA61",
	"TestConfig": "AA62",

	"ConnectionParams":        "CCC1",
	"ConnectionReqConnParams": "CCC2",
	"ConnectionDisconnReq":    "CCC3",

	"OADImageIdentify": "FFC1",
	"OADImageBlock":    "FFC2",

	"MPU9250_DATA_UUID":   "AA81",
	"MPU9250_CONFIG_UUID": "AA82",
	"MPU9250_PERIOD_UUID": "AA83",

	"LUXOMETER_CONFIG_UUID": "aa72",
	"LUXOMETER_DATA_UUID":   "aa71",
	"LUXOMETER_PERIOD_UUID": "aa73",

	"DEVICE_INFORMATION_UUID": "180A",
	"SYSTEM_ID_UUID":          "2A23",
	"MODEL_NUMBER_UUID":       "2A24",
	"SERIAL_NUMBER_UUID":      "2A25",
	"FIRMWARE_REVISION_UUID":  "2A26",
	"HARDWARE_REVISION_UUID":  "2A27",
	"SOFTWARE_REVISION_UUID":  "2A28",
	"MANUFACTURER_NAME_UUID":  "2A29",
}

//SensorTagDataEvent contains SensorTagSpecific data structure
type SensorTagDataEvent struct {
	Device     *device.Device1
	SensorType string

	AmbientTempValue interface{}
	AmbientTempUnit  string

	ObjectTempValue interface{}
	ObjectTempUnit  string

	SensorID string

	BarometericPressureValue interface{}
	BarometericPressureUnit  string

	BarometericTempValue interface{}
	BarometericTempUnit  string

	HumidityValue interface{}
	HumidityUnit  string

	HumidityTempValue interface{}
	HumidityTempUnit  string

	MpuGyroscopeValue interface{}
	MpuGyroscopeUnit  string

	MpuAccelerometerValue interface{}
	MpuAccelerometerUnit  string

	MpuMagnetometerValue interface{}
	MpuMagnetometerUnit  string

	LuxometerValue interface{}
	LuxometerUnit  string

	FirmwareVersion string
	HardwareVersion string
	Manufacturer    string
	Model           string
}

//Period =[Input*10]ms,(lowerlimit 300 ms, max 2500ms),default 1000 ms
const (
	TemperaturePeriodHigh   = 0x32  // 500 ms,
	TemperaturePeriodMedium = 0x64  // 1000 ms,
	TemperaturePeriodLow    = 0x128 // 2000 ms,
)

func getUUID(name string) (string, error) {
	if sensorTagUUIDs[name] == "" {
		return "", fmt.Errorf("Not found %s", name)
	}
	return fmt.Sprintf("F000%s-0451-4000-B000-000000000000", sensorTagUUIDs[name]), nil
}

func getDeviceInfoUUID(name string) string {
	if sensorTagUUIDs[name] == "" {
		panic("Not found " + name)
	}
	return "0000" + sensorTagUUIDs[name] + "-0000-1000-8000-00805F9B34FB"
}

//retryCall n. times, sleep millis, callback
func retryCall(times int, sleep int64, fn func() (interface{}, error)) (intf interface{}, err error) {
	for i := 0; i < times; i++ {
		intf, err = fn()
		if err == nil {
			return intf, nil
		}
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}
	return nil, err
}

//NewSensorTag creates a new sensortag instance
func NewSensorTag(d *device.Device1) (*SensorTag, error) {

	s := new(SensorTag)

	s.dataChannel = make(chan *SensorTagDataEvent)

	var connect = func(dev *device.Device1) error {
		if !dev.Properties.Connected {
			err := dev.Connect()
			if err != nil {
				return err
			}
		}
		return nil
	}

	propsChannel, err := d.WatchProperties()
	if err != nil {
		return nil, err
	}

	go func() {
		for prop := range propsChannel {
			if prop.Name == "Connected" {
				val := prop.Value.(bool)
				if val == true {
					if dataChannel != nil {
						close(dataChannel)
						break
					}
				}
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	err = connect(d)
	if err != nil {
		log.Warningf("SensorTag connection failed: %s", err)
		return nil, err
	}

	s.Device1 = d

	//initiating things for temperature sensor...(getting config,data,period characteristics...).....

	temp, err := newTemperatureSensor(s)
	if err != nil {
		return nil, err
	}
	s.Temperature = *temp

	//initiating things for humidity sensor...(getting config,data,period characteristics...).....

	humid, err := newHumiditySensor(s)
	if err != nil {
		return nil, err
	}
	s.Humidity = *humid

	//initiating things for AC,MG,GY sensor...(getting config,data,period characteristics...).....

	mpu, err := newMpuSensor(s)
	if err != nil {
		return nil, err
	}
	s.Mpu = *mpu

	//initiating things barometric sensor...(getting config,data,period characteristics...).....

	barometric, err := newBarometricSensor(s)
	if err != nil {
		return nil, err
	}
	s.Barometric = *barometric

	//initiating things luxometer sensor...(getting config,data,period characteristics...).....

	luxometer, err := newLuxometerSensor(s)
	if err != nil {
		return nil, err
	}
	s.Luxometer = *luxometer

	//initiating things for reading device info of  sensorTag...(getting firmware,hardware,manufacturer,model char...).....

	devInformation, err := newDeviceInfo(s)
	if err != nil {
		return nil, err
	}
	s.DeviceInfo = devInformation

	return s, nil

}

//SensorTag a SensorTag object representation
type SensorTag struct {
	*device.Device1
	dataChannel chan *SensorTagDataEvent
	Temperature TemperatureSensor
	Humidity    HumiditySensor
	Mpu         MpuSensor
	Barometric  BarometricSensor
	Luxometer   LuxometerSensor
	DeviceInfo  SensorTagDeviceInfo
}

func (s *SensorTag) Data() chan *SensorTagDataEvent {
	return s.dataChannel
}

//Sensor generic sensor interface
type Sensor interface {
	GetName() string
	IsEnabled() (bool, error)
	Enable() error
	Disable() error
}

func newDeviceInfo(tag *SensorTag) (SensorTagDeviceInfo, error) {

	dev := tag.Device1

	DeviceFirmwareUUID := getDeviceInfoUUID("FIRMWARE_REVISION_UUID")
	DeviceHardwareUUID := getDeviceInfoUUID("HARDWARE_REVISION_UUID")
	DeviceManufacturerUUID := getDeviceInfoUUID("MANUFACTURER_NAME_UUID")
	DeviceModelUUID := getDeviceInfoUUID("MODEL_NUMBER_UUID")

	var loadChars func() (SensorTagDeviceInfo, error)

	loadChars = func() (SensorTagDeviceInfo, error) {

		firmwareInfo, err := dev.GetCharByUUID(DeviceFirmwareUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if firmwareInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceFirmwareUUID characteristic " + DeviceFirmwareUUID)
		}

		hardwareInfo, err := dev.GetCharByUUID(DeviceHardwareUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if hardwareInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceHardwareUUID characteristic " + DeviceHardwareUUID)
		}

		manufacturerInfo, err := dev.GetCharByUUID(DeviceManufacturerUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if manufacturerInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceManufacturerUUID characteristic " + DeviceManufacturerUUID)
		}

		modelInfo, err := dev.GetCharByUUID(DeviceModelUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if modelInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceModelUUID characteristic " + DeviceModelUUID)
		}

		return SensorTagDeviceInfo{tag, modelInfo, manufacturerInfo, hardwareInfo, firmwareInfo}, err
	}

	return loadChars()
}

//SensorTagDeviceInfo sensorTag structure
type SensorTagDeviceInfo struct {
	tag              *SensorTag
	firmwareInfo     *gatt.GattCharacteristic1
	hardwareInfo     *gatt.GattCharacteristic1
	manufacturerInfo *gatt.GattCharacteristic1
	modelInfo        *gatt.GattCharacteristic1
}

//Read device info from sensorTag
func (s *SensorTagDeviceInfo) Read() (*SensorTagDataEvent, error) {

	options1 := getOptions()
	fw, err := s.firmwareInfo.ReadValue(options1)
	if err != nil {
		return nil, err
	}
	options2 := getOptions()
	hw, err := s.hardwareInfo.ReadValue(options2)
	if err != nil {
		return nil, err
	}
	options3 := getOptions()
	manufacturer, err := s.manufacturerInfo.ReadValue(options3)
	if err != nil {
		return nil, err
	}
	options4 := getOptions()
	model, err := s.modelInfo.ReadValue(options4)
	if err != nil {
		return nil, err
	}
	dataEvent := SensorTagDataEvent{
		FirmwareVersion: string(fw),
		HardwareVersion: string(hw),
		Manufacturer:    string(manufacturer),
		Model:           string(model),
	}
	return &dataEvent, err
}
