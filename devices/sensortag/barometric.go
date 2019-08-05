package sensortag

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//getting config,data,period characteristics for BAROMETRIC sensor
func newBarometricSensor(tag *SensorTag) (*BarometricSensor, error) {

	dev := tag.Device1

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
	options := getOptions()
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
	options := getOptions()
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//IsEnabled check if BarometricSensor measurements are enabled
func (s *BarometricSensor) IsEnabled() (bool, error) {
	options := getOptions()

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

	options := getOptions()
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

	err := s.Enable()
	if err != nil {
		return err
	}

	propsChanged, err := s.data.WatchProperties()
	if err != nil {
		return err
	}

	go func() {
		for prop := range propsChanged {

			if prop == nil {
				return
			}

			if prop.Name != "Value" {
				return
			}

			b1 := prop.Value.([]byte)

			barometer := binary.LittleEndian.Uint32(b1[2:])
			barometericPressureValue := calcBarometricPressure(uint32(barometer))

			barometerTemperature := binary.LittleEndian.Uint32(b1[0:4])
			barometerTempValue := calcBarometricTemperature(uint32(barometerTemperature))

			dataEvent := SensorTagDataEvent{
				Device:                   s.tag.Device1,
				SensorType:               "pressure",
				BarometericPressureValue: barometericPressureValue,
				BarometericPressureUnit:  "hPa",
				BarometericTempValue:     barometerTempValue,
				BarometericTempUnit:      "C",
				SensorID:                 macAddress,
			}

			s.tag.Data() <- &dataEvent

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
