package sensortag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//getting config,data,period characteristics for BAROMETRIC sensor
func newBarometricSensor(tag *SensorTag) (*BarometricSensor, error) {

	dev := tag.Device

	BarometerConfigUUID, err := getUUID("BarometerConfig")
	if err != nil {
		return nil, err
	}
	BarometerDataUUID, err := getUUID("BarometerData")
	if err != nil {
		return nil, err
	}
	BarometerPeriodUUID, err := getUUID("BarometerPeriod")
	if err != nil {
		return nil, err
	}

	i, err := retryCall(DefaultRetry, DefaultRetryWait, func() (interface{}, error) {

		cfg, err := dev.GetCharByUUID(BarometerConfigUUID)
		if err != nil {
			return nil, err
		}

		data, err := dev.GetCharByUUID(BarometerDataUUID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("Cannot find BarometerData characteristic " + BarometerDataUUID)
		}

		period, err := dev.GetCharByUUID(BarometerPeriodUUID)
		if err != nil {
			return nil, err
		}
		if period == nil {
			return nil, errors.New("Cannot find BarometerPeriod characteristic " + BarometerPeriodUUID)
		}

		return &BarometricSensor{tag, cfg, data, period}, err
	})

	return i.(*BarometricSensor), err
}

//BarometricSensor structure
type BarometricSensor struct {
	tag    *SensorTag
	cfg    *gatt.GattCharacteristic1
	data   *gatt.GattCharacteristic1
	period *gatt.GattCharacteristic1
}

//GetName return the sensor name
func (s BarometricSensor) GetName() string {
	return "pressure"
}

//Enable barometric measurements
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

//Disable barometric measurements
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

//IsEnabled check if BarometricSensor measurements are enabled
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

//IsNotifying check if BarometricSensor sensors are Notifying
func (s *BarometricSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}

	return n.Value().(bool), nil
}

//Read value from the BarometricSensor sensors
func (s *BarometricSensor) Read() (float64, error) {

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)

	if err != nil {
		return 0, err
	}

	amb := binary.LittleEndian.Uint16(b[2:])
	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//StartNotify enable BarometricSensorDataChannel
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

				if event1.Body[0] != gatt.GattCharacteristic1Interface {

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

				dataEvent := SensorTagDataEvent{

					Device:                   s.tag.Device,
					SensorType:               "pressure",
					BarometericPressureValue: barometericPressureValue,
					BarometericPressureUnit:  "hPa",
					BarometericTempValue:     barometerTempValue,
					BarometericTempUnit:      "C",
					SensorID:                 macAddress,
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

//StopNotify disable Barometric Sensor DataChannel
func (s *BarometricSensor) StopNotify() error {

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
