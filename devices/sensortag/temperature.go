package sensortag

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//.....getting config,data,period characteristics for TEMPERATURE sensor............
func newTemperatureSensor(tag *SensorTag) (*TemperatureSensor, error) {

	dev := tag.Device1

	TemperatureConfigUUID, err := getUUID("TemperatureConfig")
	if err != nil {
		return nil, err
	}
	TemperatureDataUUID, err := getUUID("TemperatureData")
	if err != nil {
		return nil, err
	}
	TemperaturePeriodUUID, err := getUUID("TemperaturePeriod")
	if err != nil {
		return nil, err
	}

	i, err := retryCall(DefaultRetry, DefaultRetryWait, func() (interface{}, error) {

		cfg, err := dev.GetCharByUUID(TemperatureConfigUUID)
		if err != nil {
			return nil, err
		}

		data, err := dev.GetCharByUUID(TemperatureDataUUID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, fmt.Errorf("Cannot find TemperatureData characteristic %s", TemperatureDataUUID)
		}

		period, err := dev.GetCharByUUID(TemperaturePeriodUUID)
		if err != nil {
			return nil, err
		}
		if period == nil {
			return nil, fmt.Errorf("Cannot find TemperaturePeriod characteristic %s", TemperaturePeriodUUID)
		}

		return &TemperatureSensor{tag, cfg, data, period}, err
	})

	if err != nil {
		return nil, err
	}

	return i.(*TemperatureSensor), nil
}

//TemperatureSensor the temperature sensor structure
type TemperatureSensor struct {
	tag    *SensorTag
	cfg    *gatt.GattCharacteristic1
	data   *gatt.GattCharacteristic1
	period *gatt.GattCharacteristic1
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
	options := getOptions()
	err = s.cfg.WriteValue([]byte{1}, options)
	return err
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
	options := getOptions()
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//IsEnabled check if measurements are enabled
func (s *TemperatureSensor) IsEnabled() (bool, error) {

	options := make(map[string]interface{})
	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		// TODO investigate why it report EOF only on that
		return false, nil
	}

	return (enabled == 1), nil
}

//IsNotifying check if notyfing
func (s *TemperatureSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetNotifying()
	if err != nil {
		return false, err
	}
	return n, nil
}

//Read value from the sensor
func (s *TemperatureSensor) Read() (float64, error) {

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

//StartNotify enable temperature DataChannel
func (s *TemperatureSensor) StartNotify() error {

	err := s.Enable()
	if err != nil {
		return err
	}

	dataChannel, err := s.data.WatchProperties()
	if err != nil {
		return err
	}

	go func() {
		for prop := range dataChannel {

			if prop == nil {
				return
			}

			if prop.Name != "Value" {
				return
			}

			b := prop.Value.([]byte)
			amb := binary.LittleEndian.Uint16(b[2:])
			ambientValue := calcTmpLocal(uint16(amb))

			die := binary.LittleEndian.Uint16(b[0:2])
			dieValue := calcTmpTarget(uint16(die))
			dataEvent := SensorTagDataEvent{
				Device:           s.tag.Device1,
				SensorType:       "temperature",
				AmbientTempValue: ambientValue,
				AmbientTempUnit:  "C",
				ObjectTempValue:  dieValue,
				ObjectTempUnit:   "C",
				SensorID:         s.tag.Properties.Address,
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

//StopNotify disable temperature DataChannel
func (s *TemperatureSensor) StopNotify() error {

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
