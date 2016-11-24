package devices

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/op/go-logging"
	"github.com/tj/go-debug"
)

var logger = logging.MustGetLogger("main")
var dbgTag = debug.Debug("bluez:sensortag")

var temperatureDataChannel chan dbus.Signal

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

//Sensor generic sensor interface
type Sensor interface {
	GetName() string
	IsEnabled() (bool, error)
	Enable() error
	Disable() error
}

//TemperatureSensor the temperature sensor
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

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor
var calcTmpLocal = func(raw uint16) float64 {
	return float64(raw) / 128.0
}

/* Conversion algorithm for target temperature */
// var calcTmpTarget = func(raw uint16) float64 {
//
// 	//-- calculate target temperature [°C] -
// 	Vobj2 := float64(raw) * 0.00000015625
// 	Tdie2 := calcTmpLocal(raw) + 273.15
//
// 	const S0 = 6.4E-14 // Calibration factor
// 	const a1 = 1.75E-3
// 	const a2 = -1.678E-5
// 	const b0 = -2.94E-5
// 	const b1 = -5.7E-7
// 	const b2 = 4.63E-9
// 	const c2 = 13.4
// 	const Tref = 298.15
//
// 	S := S0 * (1 + a1*(Tdie2-Tref) + a2*math.Pow((Tdie2-Tref), 2))
// 	Vos := b0 + b1*(Tdie2-Tref) + b2*math.Pow((Tdie2-Tref), 2)
// 	fObj := (Vobj2 - Vos) + c2*math.Pow((Vobj2-Vos), 2)
//
// 	tObj := math.Pow(math.Pow(Tdie2, 4)+(fObj/S), .25)
// 	tObj = (tObj - 273.15)
//
// 	return tObj
// }

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

	// die := binary.LittleEndian.Uint16(b[0:2])
	amb := binary.LittleEndian.Uint16(b[2:])

	// dieValue := calcTmpTarget(uint16(die))
	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//StartNotify enable temperatureDataChannel
// func (s *TemperatureSensor) StartNotify(fn func(temperature float64)) error {
func (s *TemperatureSensor) StartNotify() error {

	dbgTag("Enabling temperatureDataChannel")

	err := s.Enable()
	if err != nil {
		return err
	}

	temperatureDataChannel, err := s.data.Register()
	// _, err = s.data.Register()
	if err != nil {
		return err
	}

	go func() {
		for event := range temperatureDataChannel {

			if event == nil {
				return
			}

			// dbgTag("Got update %v", event)

			switch event.Body[0].(type) {
			case dbus.ObjectPath:
				// dbgTag("Received body type does not match: [0] %v -> [1] %v", event.Body[0], event.Body[1])
				continue
			case string:
				// dbgTag("body type match")
			}

			if event.Body[0] != bluez.GattCharacteristic1Interface {
				// dbgTag("Skip interface %s", event.Body[0])
				continue
			}

			props := event.Body[1].(map[string]dbus.Variant)

			if _, ok := props["Value"]; !ok {
				// dbgTag("Cannot read Value property %v", props)
				continue
			}

			b := props["Value"].Value().([]byte)

			dbgTag("Read data: %v", b)

			amb := binary.LittleEndian.Uint16(b[2:])
			ambientValue := calcTmpLocal(uint16(amb))

			// die := binary.LittleEndian.Uint16(b[0:2])
			// dieValue := calcTmpTarget(uint16(die))

			dbgTag("Got data %v", ambientValue)

			dataEvent := api.DataEvent{
				Device: s.tag.Device,
				Sensor: "temperature",
				Value:  ambientValue,
				Unit:   "C",
			}
			s.tag.Device.Emit("data", dataEvent)

		}
	}()

	return s.data.StartNotify()
}

//StopNotify disable temperatureDataChannel
func (s *TemperatureSensor) StopNotify() error {

	dbgTag("Disabling temperatureDataChannel")

	err := s.Disable()
	if err != nil {
		return err
	}

	if temperatureDataChannel != nil {
		close(temperatureDataChannel)
	}

	return s.data.StopNotify()
}

// NewSensorTag creates a new sensortag instance
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

				if temperatureDataChannel != nil {
					close(temperatureDataChannel)
				}

			}
		}
	}))

	err := connect(d)
	if err != nil {
		logger.Warning("Connection failed %s", err)
		return nil, err
	}

	s.Device = d

	temp, err := newTemperatureSensor(s)
	if err != nil {
		return nil, err
	}
	s.Temperature = temp

	return s, nil
}

//SensorTag a SensorTag object representation
type SensorTag struct {
	*api.Device
	Temperature TemperatureSensor
}
