package devices

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/tj/go-debug"
)

var dbgTag = debug.Debug("bluetooth:sensortag")

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
	Device     *api.Device
	SensorType string

	AmbientTempValue interface{}
	AmbientTempUnit  string

	ObjectTempValue interface{}
	ObjectTempUnit  string

	SensorId string

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

func getUUID(name string) string {
	if sensorTagUUIDs[name] == "" {
		panic("Not found " + name)
	}
	return "F000" + sensorTagUUIDs[name] + "-0451-4000-B000-000000000000"
}

func getDeviceInfoUUID(name string) string {
	if sensorTagUUIDs[name] == "" {
		panic("Not found " + name)
	}
	return "0000" + sensorTagUUIDs[name] + "-0000-1000-8000-00805F9B34FB"
}

//....getting config,data,period characteristics for Humidity sensor...........................

func newHumiditySensor(tag *SensorTag) (HumiditySensor, error) {

	dev := tag.Device

	HumidityConfigUUID := getUUID("HumidityConfig")
	HumidityDataUUID := getUUID("HumidityData")
	HumidityPeriodUUID := getUUID("HumidityPeriod")

	retry := 3
	tries := 0
	var loadChars func() (HumiditySensor, error)

	loadChars = func() (HumiditySensor, error) {

		dbgTag("Load humid cfg")
		cfg, err := dev.GetCharByUUID(HumidityConfigUUID)

		if err != nil {
			return HumiditySensor{}, err
		}

		if cfg == nil {

			if tries == retry {
				return HumiditySensor{}, errors.New("Cannot find HumidityConfig characteristic " + HumidityConfigUUID)
			}

			tries++
			time.Sleep(time.Second * time.Duration(5*tries))
			dbgTag("Char not found, try to reload")

			return loadChars()
		}

		dbgTag("Load humid data")
		data, err := dev.GetCharByUUID(HumidityDataUUID)
		if err != nil {
			return HumiditySensor{}, err
		}
		if data == nil {
			return HumiditySensor{}, errors.New("Cannot find HumidityData characteristic " + HumidityDataUUID)
		}

		dbgTag("Load temp period")
		period, err := dev.GetCharByUUID(HumidityPeriodUUID)
		if err != nil {
			return HumiditySensor{}, err
		}
		if period == nil {
			return HumiditySensor{}, errors.New("Cannot find HumidityPeriod characteristic " + HumidityPeriodUUID)
		}

		return HumiditySensor{tag, cfg, data, period}, err
	}

	return loadChars()
}

//......humidity sensor structure....

type HumiditySensor struct {
	tag    *SensorTag
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// ........GetName return the sensor name..............

func (s HumiditySensor) GetName() string {
	return "humidity"
}

//.......Enable humidity measurements.................

func (s *HumiditySensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{1}, options)
	if err != nil {
		return err
	}
	return nil
}

//............Disable humidity measurements................

func (s *HumiditySensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//.........IsEnabled check if humidity measurements are enabled...........

func (s *HumiditySensor) IsEnabled() (bool, error) {
	options := make(map[string]dbus.Variant)

	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	return (enabled == 1), nil
}

//.......IsNotifying check if humidity sensor is notyfing................

func (s *HumiditySensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.(bool), nil
}

//........Read value from the humidity sensor........................

func (s *HumiditySensor) Read() (float64, error) {

	dbgTag("Reading humidity sensor")

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)

	dbgTag("Read data: %v", b)

	if err != nil {
		return 0, err
	}

	//die := binary.LittleEndian.Uint16(b[0:2])
	amb := binary.LittleEndian.Uint16(b[2:])

	//dieValue := calcTmpTarget(uint16(die))
	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//..........StartNotify enable DataChannel for humidity................

func (s *HumiditySensor) StartNotify(macAddress string) error {

	d := s.tag.Device
	serv, err1 := d.GetAllServicesAndUUID()

	if err1 != nil {

	}
	var uuidAndService string
	serviceArrLength := len(serv)
	for i := 0; i < serviceArrLength; i++ {

		val := strings.Split(serv[i], ":")

		if val[0] == "F000AA22-0451-4000-B000-000000000000" {
			uuidAndService = val[1]
		}
	}
	dbgTag("Enabling dataChannel for humidity")

	err := s.Enable()
	if err != nil {
		return err
	}

	dataChannel, err := s.data.Register()

	if err != nil {
		return err
	}

	go func() {
		for event1 := range dataChannel {

			if event1 == nil {
				return
			}

			if strings.Contains(fmt.Sprint(event1.Path), uuidAndService) {

				switch event1.Body[0].(type) {

				case dbus.ObjectPath:
					continue

				case string:
				}

				if event1.Body[0] != bluez.GattCharacteristic1Interface {

					continue
				}

				props1 := event1.Body[1].(map[string]dbus.Variant)

				if _, ok := props1["Value"]; !ok {
					continue
				}

				b1 := props1["Value"].Value().([]byte)
				dbgTag("Read data: %v", b1)

				humid := binary.LittleEndian.Uint16(b1[2:])

				humidityValue := calcHumidLocal(uint16(humid))

				temperature := binary.LittleEndian.Uint16(b1[0:2])
				tempValue := calcTmpFromHumidSensor(uint16(temperature))
				dbgTag("Got data %v", humidityValue)
				dataEvent := SensorTagDataEvent{

					Device:            s.tag.Device,
					SensorType:        "humidity",
					HumidityValue:     humidityValue,
					HumidityUnit:      "%RH",
					HumidityTempValue: tempValue,
					HumidityTempUnit:  "C",
					SensorId:          macAddress,
				}
				s.tag.Device.Emit("data", dataEvent)
			}
		}
	}()

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if !n {
		return s.data.StartNotify()
	}
	return nil
}

//............StopNotify disable DataChannel for humidity sensor............

func (s *HumiditySensor) StopNotify() error {

	dbgTag("Disabling dataChannel")

	err := s.Disable()
	if err != nil {
		return err
	}

	if dataChannel != nil {
		close(dataChannel)
	}

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if n {
		return s.data.StopNotify()
	}
	return nil
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor

var calcHumidLocal = func(raw uint16) float64 {

	//.....................humidity calibiration.......................

	return float64(raw) * 100 / 65536.0
}

var calcTmpFromHumidSensor = func(raw uint16) float64 {

	//.......TEMPERATURE calibiration for data comming from humidity sensor..............

	return -40 + ((165 * float64(raw)) / 65536.0)

}

//..................................end of humidity sensor.....................................

//.................................MPU SENSORS (Accelerometer,magnetometer,gyroscope)..........

//....getting config,data,period characteristics for Humidity sensor...........................

func newMpuSensor(tag *SensorTag) (MpuSensor, error) {

	dev := tag.Device

	//...........accelerometer,magnetometer,gyroscope..........

	MpuConfigUUID := getUUID("MPU9250_CONFIG_UUID")
	MpuDataUUID := getUUID("MPU9250_DATA_UUID")
	MpuPeriodUUID := getUUID("MPU9250_PERIOD_UUID")

	retry := 3
	tries := 0
	var loadChars func() (MpuSensor, error)

	loadChars = func() (MpuSensor, error) {

		dbgTag("Load mpu cfg")
		cfg, err := dev.GetCharByUUID(MpuConfigUUID)

		if err != nil {
			return MpuSensor{}, err
		}

		if cfg == nil {

			if tries == retry {
				return MpuSensor{}, errors.New("Cannot find MpuConfig characteristic " + MpuConfigUUID)
			}

			tries++
			time.Sleep(time.Second * time.Duration(5*tries))
			dbgTag("Char not found, try to reload")

			return loadChars()
		}

		dbgTag("Load mpu data")
		data, err := dev.GetCharByUUID(MpuDataUUID)
		if err != nil {
			return MpuSensor{}, err
		}
		if data == nil {
			return MpuSensor{}, errors.New("Cannot find MpuData characteristic " + MpuDataUUID)
		}

		dbgTag("Load mpu period")
		period, err := dev.GetCharByUUID(MpuPeriodUUID)
		if err != nil {
			return MpuSensor{}, err
		}
		if period == nil {
			return MpuSensor{}, errors.New("Cannot find MpuPeriod characteristic " + MpuPeriodUUID)
		}

		return MpuSensor{tag, cfg, data, period}, err
	}

	return loadChars()
}

//......Mpu Sensor structure..........

type MpuSensor struct {
	tag    *SensorTag
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// ........GetName return's the sensor name..............

func (s MpuSensor) GetName() string {
	return "Ac-Mg-Gy"
}

//.......Enable mpuSensors measurements.................

func (s *MpuSensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{0x0007f, 0x0007f}, options)
	if err != nil {
		return err
	}
	return nil
}

//............Disable mpuSensors measurements................

func (s *MpuSensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//.........IsEnabled check if mpu measurements are enabled...........

func (s *MpuSensor) IsEnabled() (bool, error) {
	options := make(map[string]dbus.Variant)

	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	return (enabled == 1), nil
}

//.......IsNotifying check if mpu sensors are notifying................

func (s *MpuSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.(bool), nil
}

//........Read value from the mpu sensors........................

func (s *MpuSensor) Read() (float64, error) {

	dbgTag("Reading mpu sensor")

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)

	dbgTag("Read data: %v", b)

	if err != nil {
		return 0, err
	}
	amb := binary.LittleEndian.Uint16(b[2:])

	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//..........StartNotify enable mpuDataChannel................

func (s *MpuSensor) StartNotify(macAddress string) error {

	//log.Debug("MpuSensor tag value: ",s.tag.Device)
	d := s.tag.Device
	serv, err1 := d.GetAllServicesAndUUID()

	if err1 != nil {

	}
	var uuidAndService string
	serviceArrLength := len(serv)
	for i := 0; i < serviceArrLength; i++ {
		val := strings.Split(serv[i], ":")

		if val[0] == "F000AA82-0451-4000-B000-000000000000" {
			uuidAndService = val[1]
		}
	}
	dbgTag("Enabling mpuDataChannel")

	err := s.Enable()
	if err != nil {
		return err
	}

	dataChannel, err := s.data.Register()
	if err != nil {
		return err
	}

	go func() {
		for event1 := range dataChannel {

			if event1 == nil {
				return
			}

			if strings.Contains(fmt.Sprint(event1.Path), uuidAndService) {
				switch event1.Body[0].(type) {

				case dbus.ObjectPath:
					continue
				case string:
				}

				if event1.Body[0] != bluez.GattCharacteristic1Interface {

					continue
				}

				props1 := event1.Body[1].(map[string]dbus.Variant)

				if _, ok := props1["Value"]; !ok {

					continue
				}

				b1 := props1["Value"].Value().([]byte)
				var mpuAccelerometer string
				var mpuGyroscope string
				var mpuMagnetometer string

				//.......... calculate Gyroscope .................................

				mpuXg := binary.LittleEndian.Uint16(b1[0:2])
				mpuYg := binary.LittleEndian.Uint16(b1[2:4])
				mpuZg := binary.LittleEndian.Uint16(b1[4:6])

				mpuGyX, mpuGyY, mpuGyZ := calcMpuGyroscope(uint16(mpuXg), uint16(mpuYg), uint16(mpuZg))
				mpuGyroscope = fmt.Sprint(mpuGyX, " , ", mpuGyY, " , ", mpuGyZ)

				//.......... calculate Accelerometer .............................

				mpuXa := binary.LittleEndian.Uint16(b1[6:8])
				mpuYa := binary.LittleEndian.Uint16(b1[8:10])
				mpuZa := binary.LittleEndian.Uint16(b1[10:12])

				mpuAcX, mpuAcY, mpuAcZ := calcMpuAccelerometer(uint16(mpuXa), uint16(mpuYa), uint16(mpuZa))
				mpuAccelerometer = fmt.Sprint(mpuAcX, " , ", mpuAcY, " , ", mpuAcZ)

				//.......... calculate Magnetometer .............................

				mpuXm := binary.LittleEndian.Uint16(b1[12:14])
				mpuYm := binary.LittleEndian.Uint16(b1[14:16])
				mpuZm := binary.LittleEndian.Uint16(b1[16:18])

				mpuMgX, mpuMgY, mpuMgZ := calcMpuMagnetometer(uint16(mpuXm), uint16(mpuYm), uint16(mpuZm))
				mpuMagnetometer = fmt.Sprint(mpuMgX, " , ", mpuMgY, " , ", mpuMgZ)

				dataEvent := SensorTagDataEvent{

					Device:                s.tag.Device,
					SensorType:            "mpu",
					MpuGyroscopeValue:     mpuGyroscope,
					MpuGyroscopeUnit:      "deg/s",
					MpuAccelerometerValue: mpuAccelerometer,
					MpuAccelerometerUnit:  "G",
					MpuMagnetometerValue:  mpuMagnetometer,
					MpuMagnetometerUnit:   "uT",
					SensorId:              macAddress,
				}
				s.tag.Device.Emit("data", dataEvent)
			}
		}
	}()

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if !n {
		return s.data.StartNotify()
	}
	return nil
}

//............StopNotify disable DataChannel for mpu sensors............

func (s *MpuSensor) StopNotify() error {

	dbgTag("Disabling dataChannel")

	err := s.Disable()
	if err != nil {
		return err
	}

	if dataChannel != nil {
		close(dataChannel)
	}

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if n {
		return s.data.StopNotify()
	}
	return nil
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor

var calcMpuGyroscope = func(rawX, rawY, rawZ uint16) (float64, float64, float64) {

	Xg := float64(rawX) / 128.0
	Yg := float64(rawY) / 128.0
	Zg := float64(rawZ) / 128.0

	return Xg, Yg, Zg
}
var calcMpuAccelerometer = func(rawX, rawY, rawZ uint16) (float64, float64, float64) {

	Xg := float64(rawX) / 4096.0
	Yg := float64(rawY) / 4096.0
	Zg := float64(rawZ) / 4096.0

	return Xg, Yg, Zg
}
var calcMpuMagnetometer = func(rawX, rawY, rawZ uint16) (float64, float64, float64) {

	Xg := float64(rawX) * 4912.0 / 32768.0
	Yg := float64(rawY) * 4912.0 / 32768.0
	Zg := float64(rawZ) * 4912.0 / 32768.0

	return Xg, Yg, Zg
}

//.................................. end of MpuSensor .........................................

//.................................BAROMETRIC SENSORS..........................................

//....getting config,data,period characteristics for BAROMETRIC sensor.........................

func newBarometricSensor(tag *SensorTag) (BarometricSensor, error) {

	dev := tag.Device

	BarometerConfigUUID := getUUID("BarometerConfig")
	BarometerDataUUID := getUUID("BarometerData")
	BarometerPeriodUUID := getUUID("BarometerPeriod")

	retry := 3
	tries := 0
	var loadChars func() (BarometricSensor, error)

	loadChars = func() (BarometricSensor, error) {

		dbgTag("Load pressure cfg")
		cfg, err := dev.GetCharByUUID(BarometerConfigUUID)

		if err != nil {
			return BarometricSensor{}, err
		}

		if cfg == nil {

			if tries == retry {
				return BarometricSensor{}, errors.New("Cannot find BarometerConfig characteristic " + BarometerConfigUUID)
			}

			tries++
			time.Sleep(time.Second * time.Duration(5*tries))
			dbgTag("Char not found, try to reload")

			return loadChars()
		}

		dbgTag("Load barometer data")
		data, err := dev.GetCharByUUID(BarometerDataUUID)
		if err != nil {
			return BarometricSensor{}, err
		}
		if data == nil {
			return BarometricSensor{}, errors.New("Cannot find BarometerData characteristic " + BarometerDataUUID)
		}

		dbgTag("Load barometer period")
		period, err := dev.GetCharByUUID(BarometerPeriodUUID)
		if err != nil {
			return BarometricSensor{}, err
		}
		if period == nil {
			return BarometricSensor{}, errors.New("Cannot find BarometerPeriod characteristic " + BarometerPeriodUUID)
		}

		return BarometricSensor{tag, cfg, data, period}, err
	}

	return loadChars()
}

//......Barometric sensor structure..........

type BarometricSensor struct {
	tag    *SensorTag
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// ........GetName return the sensor name..............

func (s BarometricSensor) GetName() string {
	return "pressure"
}

//.......Enable barometric measurements.................

func (s *BarometricSensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)

	err = s.cfg.WriteValue([]byte{1}, options)
	if err != nil {
		return err
	}
	return nil
}

//............Disable barometric measurements................

func (s *BarometricSensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//.........IsEnabled check if BarometricSensor measurements are enabled...........

func (s *BarometricSensor) IsEnabled() (bool, error) {
	options := make(map[string]dbus.Variant)

	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	return (enabled == 1), nil
}

//.......IsNotifying check if BarometricSensor sensors are Notifying.........

func (s *BarometricSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.(bool), nil
}

//........Read value from the BarometricSensor sensors........................

func (s *BarometricSensor) Read() (float64, error) {

	dbgTag("Reading BarometricSensor sensor")

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)

	dbgTag("Read data: %v", b)

	if err != nil {
		return 0, err
	}

	amb := binary.LittleEndian.Uint16(b[2:])
	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//..........StartNotify enable BarometricSensorDataChannel................

func (s *BarometricSensor) StartNotify(macAddress string) error {

	d := s.tag.Device
	serv, err1 := d.GetAllServicesAndUUID()

	if err1 != nil {

	}
	var uuidAndService string
	serviceArrLength := len(serv)
	for i := 0; i < serviceArrLength; i++ {

		val := strings.Split(serv[i], ":")
		if val[0] == "F000AA42-0451-4000-B000-000000000000" {
			uuidAndService = val[1]
		}
	}
	dbgTag("Enabling BarometricSensorDataChannel")

	err := s.Enable()
	if err != nil {
		return err
	}

	dataChannel, err := s.data.Register()
	if err != nil {
		return err
	}

	go func() {
		for event1 := range dataChannel {

			if event1 == nil {
				return
			}
			if strings.Contains(fmt.Sprint(event1.Path), uuidAndService) {

				switch event1.Body[0].(type) {

				case dbus.ObjectPath:

					continue
				case string:
				}

				if event1.Body[0] != bluez.GattCharacteristic1Interface {

					continue
				}

				props1 := event1.Body[1].(map[string]dbus.Variant)

				if _, ok := props1["Value"]; !ok {

					continue
				}

				b1 := props1["Value"].Value().([]byte)

				barometer := binary.LittleEndian.Uint32(b1[2:])
				barometericPressureValue := calcBarometricPressure(uint32(barometer))

				barometerTemperature := binary.LittleEndian.Uint32(b1[0:4])
				barometerTempValue := calcBarometricTemperature(uint32(barometerTemperature))

				dbgTag("Got data %v", barometericPressureValue)

				dataEvent := SensorTagDataEvent{

					Device:                   s.tag.Device,
					SensorType:               "pressure",
					BarometericPressureValue: barometericPressureValue,
					BarometericPressureUnit:  "hPa",
					BarometericTempValue:     barometerTempValue,
					BarometericTempUnit:      "C",
					SensorId:                 macAddress,
				}
				s.tag.Device.Emit("data", dataEvent)
			}
		}
	}()

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if !n {
		return s.data.StartNotify()
	}
	return nil
}

//............StopNotify disable Barometric Sensor DataChannel............

func (s *BarometricSensor) StopNotify() error {

	dbgTag("Disabling dataChannel")

	err := s.Disable()
	if err != nil {
		return err
	}

	if dataChannel != nil {
		close(dataChannel)
	}

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if n {
		return s.data.StopNotify()
	}
	return nil
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor

var calcBarometricPressure = func(raw uint32) float64 {
	//.......barometric pressure......
	pressureMask := (int(raw) >> 8) & 0x00ffffff
	return float64(pressureMask) / 100.0
}

var calcBarometricTemperature = func(raw uint32) float64 {
	//.......TEMPERATURE calibiration data comming from barometric sensor...........
	tempMask := int(raw) & 0x00ffffff
	return float64(tempMask) / 100.0
}

//.................................. end of BarometricSensor ..............................

//........................................Temperature Sensor...............................

//............getting config,data,period characteristics for TEMPERATURE sensor............

func newTemperatureSensor(tag *SensorTag) (TemperatureSensor, error) {

	dev := tag.Device

	TemperatureConfigUUID := getUUID("TemperatureConfig")
	TemperatureDataUUID := getUUID("TemperatureData")
	TemperaturePeriodUUID := getUUID("TemperaturePeriod")

	retry := 3
	tries := 0
	var loadChars func() (TemperatureSensor, error)

	loadChars = func() (TemperatureSensor, error) {

		dbgTag("Load temp cfg")
		cfg, err := dev.GetCharByUUID(TemperatureConfigUUID)

		if err != nil {
			return TemperatureSensor{}, err
		}

		if cfg == nil {

			if tries == retry {
				return TemperatureSensor{}, errors.New("Cannot find TemperatureConfig characteristic " + TemperatureConfigUUID)
			}

			tries++
			time.Sleep(time.Second * time.Duration(5*tries))
			dbgTag("Char not found, try to reload")

			return loadChars()
		}

		dbgTag("Load temp data")
		data, err := dev.GetCharByUUID(TemperatureDataUUID)
		if err != nil {
			return TemperatureSensor{}, err
		}
		if data == nil {
			return TemperatureSensor{}, errors.New("Cannot find TemperatureData characteristic " + TemperatureDataUUID)
		}

		dbgTag("Load temp period")
		period, err := dev.GetCharByUUID(TemperaturePeriodUUID)
		if err != nil {
			return TemperatureSensor{}, err
		}
		if period == nil {
			return TemperatureSensor{}, errors.New("Cannot find TemperaturePeriod characteristic " + TemperaturePeriodUUID)
		}

		return TemperatureSensor{tag, cfg, data, period}, err
	}

	return loadChars()
}

//TemperatureSensor the temperature sensor structure

type TemperatureSensor struct {
	tag    *SensorTag
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// GetName return the sensor name
func (s TemperatureSensor) GetName() string {
	return "temperature"
}

//Enable measurements
func (s *TemperatureSensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{1}, options)

	if err != nil {

		return err
	}
	return nil
}

//Disable measurements
func (s *TemperatureSensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//IsEnabled check if measurements are enabled
func (s *TemperatureSensor) IsEnabled() (bool, error) {
	options := make(map[string]dbus.Variant)

	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	return (enabled == 1), nil
}

//IsNotifying check if notyfing
func (s *TemperatureSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.(bool), nil
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor
var calcTmpLocal = func(raw uint16) float64 {

	//.......ambient temperature calberation..............

	return float64(raw) / 128.0

}

/* Conversion algorithm for target temperature */
var calcTmpTarget = func(raw uint16) float64 {

	//.........object temperature caliberation...........

	return float64(raw) / 128.0
}

//Read value from the sensor
func (s *TemperatureSensor) Read() (float64, error) {

	dbgTag("Reading temperature sensor")

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)

	dbgTag("Read data: %v", b)

	if err != nil {
		return 0, err
	}

	amb := binary.LittleEndian.Uint16(b[2:])
	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//StartNotify enable temperature DataChannel

func (s *TemperatureSensor) StartNotify(macAddress string) error {

	d := s.tag.Device
	serv, err1 := d.GetAllServicesAndUUID()

	if err1 != nil {

	}
	var uuidAndService string
	serviceArrLength := len(serv)
	for i := 0; i < serviceArrLength; i++ {

		val := strings.Split(serv[i], ":")

		if val[0] == "F000AA01-0451-4000-B000-000000000000" {
			uuidAndService = val[1]
		}
	}

	dbgTag("Enabling DataChannel")

	err := s.Enable()
	if err != nil {
		return err
	}

	dataChannel, err := s.data.Register()
	if err != nil {
		return err
	}

	go func() {
		for event := range dataChannel {

			if event == nil {
				return
			}
			if strings.Contains(fmt.Sprint(event.Path), uuidAndService) {

				switch event.Body[0].(type) {

				case dbus.ObjectPath:
					continue

				case string:
				}

				if event.Body[0] != bluez.GattCharacteristic1Interface {
					continue
				}

				props := event.Body[1].(map[string]dbus.Variant)

				if _, ok := props["Value"]; !ok {
					continue
				}

				b := props["Value"].Value().([]byte)
				amb := binary.LittleEndian.Uint16(b[2:])
				ambientValue := calcTmpLocal(uint16(amb))

				die := binary.LittleEndian.Uint16(b[0:2])
				dieValue := calcTmpTarget(uint16(die))
				dataEvent := SensorTagDataEvent{

					Device:           s.tag.Device,
					SensorType:       "temperature",
					AmbientTempValue: ambientValue,
					AmbientTempUnit:  "C",
					ObjectTempValue:  dieValue,
					ObjectTempUnit:   "C",
					SensorId:         macAddress,
				}
				s.tag.Device.Emit("data", dataEvent)
			}
		}
	}()

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if !n {
		return s.data.StartNotify()
	}
	return nil
}

//StopNotify disable temperature DataChannel
func (s *TemperatureSensor) StopNotify() error {

	dbgTag("Disabling temperature DataChannel")

	err := s.Disable()
	if err != nil {
		return err
	}

	if dataChannel != nil {
		close(dataChannel)
	}

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if n {
		return s.data.StopNotify()
	}
	return nil
}

//.................................end of temperature sensor...................................

//.................................Luxometer Sensor............................................

//....getting config,data,period characteristics for luxometer sensor.........................

func newLuxometerSensor(tag *SensorTag) (LuxometerSensor, error) {

	dev := tag.Device

	LuxometerConfigUUID := getUUID("LUXOMETER_CONFIG_UUID")
	LuxometerDataUUID := getUUID("LUXOMETER_DATA_UUID")
	LuxometerPeriodUUID := getUUID("LUXOMETER_PERIOD_UUID")

	retry := 3
	tries := 0
	var loadChars func() (LuxometerSensor, error)

	loadChars = func() (LuxometerSensor, error) {

		dbgTag("Load luxometer cfg")
		cfg, err := dev.GetCharByUUID(LuxometerConfigUUID)

		if err != nil {
			return LuxometerSensor{}, err
		}

		if cfg == nil {

			if tries == retry {
				return LuxometerSensor{}, errors.New("Cannot find LuxometerConfigUUID  characteristic " + LuxometerConfigUUID)
			}

			tries++
			time.Sleep(time.Second * time.Duration(5*tries))
			dbgTag("Char not found, try to reload")

			return loadChars()
		}

		dbgTag("Load luxometer data")
		data, err := dev.GetCharByUUID(LuxometerDataUUID)
		if err != nil {
			return LuxometerSensor{}, err
		}
		if data == nil {
			return LuxometerSensor{}, errors.New("Cannot find LuxometerDataUUID  characteristic " + LuxometerDataUUID)
		}

		dbgTag("Load luxometer period")
		period, err := dev.GetCharByUUID(LuxometerPeriodUUID)
		if err != nil {
			return LuxometerSensor{}, err
		}
		if period == nil {
			return LuxometerSensor{}, errors.New("Cannot find LuxometerPeriodUUID  characteristic " + LuxometerPeriodUUID)
		}

		return LuxometerSensor{tag, cfg, data, period}, err
	}

	return loadChars()
}

//......Luxometer sensor structure..........

type LuxometerSensor struct {
	tag    *SensorTag
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// ........GetName return the sensor name..............

func (s LuxometerSensor) GetName() string {
	return "luxometer"
}

//.......Enable LuxometerSensor measurements.................

func (s *LuxometerSensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{1}, options)
	if err != nil {
		return err
	}
	return nil
}

//............Disable LuxometerSensor measurements................

func (s *LuxometerSensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]dbus.Variant)
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//.........IsEnabled check if LuxometerSensor measurements are enabled.......

func (s *LuxometerSensor) IsEnabled() (bool, error) {
	options := make(map[string]dbus.Variant)

	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	return (enabled == 1), nil
}

//.......IsNotifying check if LuxometerSensor sensors are Notifying.........

func (s *LuxometerSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.(bool), nil
}

//........Read value from the LuxometerSensor sensors........................

func (s *LuxometerSensor) Read() (float64, error) {

	dbgTag("Reading LuxometerSensor sensor")

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)

	dbgTag("Read data: %v", b)

	if err != nil {
		return 0, err
	}

	amb := binary.LittleEndian.Uint16(b[2:])
	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//..........StartNotify enable LuxometerSensorDataChannel................

func (s *LuxometerSensor) StartNotify(macAddress string) error {

	d := s.tag.Device
	serv, err1 := d.GetAllServicesAndUUID()

	if err1 != nil {

	}
	var uuidAndService string
	serviceArrLength := len(serv)
	for i := 0; i < serviceArrLength; i++ {

		val := strings.Split(serv[i], ":")

		if val[0] == "F000AA71-0451-4000-B000-000000000000" {
			uuidAndService = val[1]
		}
	}
	dbgTag("Enabling LuxometerSensorDataChannel")

	err := s.Enable()
	if err != nil {
		return err
	}

	dataChannel, err := s.data.Register()
	if err != nil {
		return err
	}

	go func() {
		for event1 := range dataChannel {

			if event1 == nil {
				return
			}
			if strings.Contains(fmt.Sprint(event1.Path), uuidAndService) {

				switch event1.Body[0].(type) {

				case dbus.ObjectPath:
					continue

				case string:
				}

				if event1.Body[0] != bluez.GattCharacteristic1Interface {

					continue
				}

				props1 := event1.Body[1].(map[string]dbus.Variant)

				if _, ok := props1["Value"]; !ok {

					continue
				}

				b1 := props1["Value"].Value().([]byte)
				luxometer := binary.LittleEndian.Uint16(b1[0:])
				luxometerValue := calcLuxometer(uint16(luxometer))

				dataEvent := SensorTagDataEvent{

					Device:         s.tag.Device,
					SensorType:     "luxometer",
					LuxometerValue: luxometerValue,
					LuxometerUnit:  "candela",
					SensorId:       macAddress,
				}
				s.tag.Device.Emit("data", dataEvent)
			}
		}
	}()

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if !n {
		return s.data.StartNotify()
	}
	return nil
}

//............StopNotify disable Luxometer Sensor DataChannel............

func (s *LuxometerSensor) StopNotify() error {

	dbgTag("Disabling dataChannel")

	err := s.Disable()
	if err != nil {
		return err
	}

	if dataChannel != nil {
		close(dataChannel)
	}

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if n {
		return s.data.StopNotify()
	}
	return nil
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor

var calcLuxometer = func(raw uint16) float64 {

	exponent := (int(raw) & 0xF000) >> 12
	mantissa := (int(raw) & 0x0FFF)
	exp := float64(exponent)
	man := float64(mantissa)
	flLux := man * math.Pow(2, exp) / 100
	return float64(flLux)
}

//.................................. end of LuxometerSensor .....................................

//.....NewSensorTag creates a new sensortag instance.....

func NewSensorTag(d *api.Device) (*SensorTag, error) {

	s := new(SensorTag)

	var connect = func(dev *api.Device) error {
		if !dev.IsConnected() {
			err := dev.Connect()
			if err != nil {
				return err
			}
		}
		return nil
	}

	d.On("changed", emitter.NewCallback(func(ev emitter.Event) {
		changed := ev.GetData().(api.PropertyChangedEvent)
		if changed.Field == "Connected" {
			conn := changed.Value.(bool)
			if !conn {

				dbgTag("Device disconnected")

				// TODO clean up properly

				if dataChannel != nil {
					close(dataChannel)
				}

			}
		}
	}))

	err := connect(d)
	if err != nil {
		log.Warning("SensorTag connection failed: %v", err)
		return nil, err
	}

	s.Device = d

	//initiating things for temperature sensor...(getting config,data,period characteristics...).....

	temp, err := newTemperatureSensor(s)
	if err != nil {
		return nil, err
	}
	s.Temperature = temp

	//initiating things for humidity sensor...(getting config,data,period characteristics...).....

	humid, err := newHumiditySensor(s)
	if err != nil {
		return nil, err
	}
	s.Humidity = humid

	//initiating things for AC,MG,GY sensor...(getting config,data,period characteristics...).....

	mpu, err := newMpuSensor(s)
	if err != nil {
		return nil, err
	}
	s.Mpu = mpu

	//initiating things barometric sensor...(getting config,data,period characteristics...).....

	barometric, err := newBarometricSensor(s)
	if err != nil {
		return nil, err
	}
	s.Barometric = barometric

	//initiating things luxometer sensor...(getting config,data,period characteristics...).....

	luxometer, err := newLuxometerSensor(s)
	if err != nil {
		return nil, err
	}
	s.Luxometer = luxometer

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
	*api.Device
	Temperature TemperatureSensor
	Humidity    HumiditySensor
	Mpu         MpuSensor
	Barometric  BarometricSensor
	Luxometer   LuxometerSensor
	DeviceInfo  SensorTagDeviceInfo
}

//Sensor generic sensor interface
type Sensor interface {
	GetName() string
	IsEnabled() (bool, error)
	Enable() error
	Disable() error
}

func newDeviceInfo(tag *SensorTag) (SensorTagDeviceInfo, error) {

	dev := tag.Device

	DeviceFirmwareUUID := getDeviceInfoUUID("FIRMWARE_REVISION_UUID")
	DeviceHardwareUUID := getDeviceInfoUUID("HARDWARE_REVISION_UUID")
	DeviceManufacturerUUID := getDeviceInfoUUID("MANUFACTURER_NAME_UUID")
	DeviceModelUUID := getDeviceInfoUUID("MODEL_NUMBER_UUID")

	var loadChars func() (SensorTagDeviceInfo, error)

	loadChars = func() (SensorTagDeviceInfo, error) {

		dbgTag("Load device DeviceFirmwareUUID")

		firmwareInfo, err := dev.GetCharByUUID(DeviceFirmwareUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if firmwareInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceFirmwareUUID characteristic " + DeviceFirmwareUUID)
		}

		dbgTag("Load device DeviceHardwareUUID")

		hardwareInfo, err := dev.GetCharByUUID(DeviceHardwareUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if hardwareInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceHardwareUUID characteristic " + DeviceHardwareUUID)
		}

		dbgTag("Load device DeviceManufacturerUUID")

		manufacturerInfo, err := dev.GetCharByUUID(DeviceManufacturerUUID)
		if err != nil {
			return SensorTagDeviceInfo{}, err
		}
		if manufacturerInfo == nil {
			return SensorTagDeviceInfo{}, errors.New("Cannot find DeviceManufacturerUUID characteristic " + DeviceManufacturerUUID)
		}

		dbgTag("Load device DeviceModelUUID")

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
	firmwareInfo     *profile.GattCharacteristic1
	hardwareInfo     *profile.GattCharacteristic1
	manufacturerInfo *profile.GattCharacteristic1
	modelInfo        *profile.GattCharacteristic1
}

//Read device info from sensorTag
func (s *SensorTagDeviceInfo) Read() (SensorTagDataEvent, error) {

	options1 := make(map[string]dbus.Variant)
	fw, err := s.firmwareInfo.ReadValue(options1)

	options2 := make(map[string]dbus.Variant)
	hw, err := s.hardwareInfo.ReadValue(options2)

	options3 := make(map[string]dbus.Variant)
	manufacturer, err := s.manufacturerInfo.ReadValue(options3)

	options4 := make(map[string]dbus.Variant)
	model, err := s.modelInfo.ReadValue(options4)

	dataEvent := SensorTagDataEvent{

		FirmwareVersion: string(fw),
		HardwareVersion: string(hw),
		Manufacturer:    string(manufacturer),
		Model:           string(model),
	}
	return dataEvent, err
}
