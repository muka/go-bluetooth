package sensortag

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//getting config,data,period characteristics for Humidity sensor
func newHumiditySensor(tag *SensorTag) (*HumiditySensor, error) {

	dev := tag.Device1
	HumidityConfigUUID, err := getUUID("HumidityConfig")
	if err != nil {
		return nil, err
	}
	HumidityDataUUID, err := getUUID("HumidityData")
	if err != nil {
		return nil, err
	}
	HumidityPeriodUUID, err := getUUID("HumidityPeriod")
	if err != nil {
		return nil, err
	}

	i, err := retryCall(DefaultRetry, DefaultRetryWait, func() (interface{}, error) {
		cfg, err := dev.GetCharByUUID(HumidityConfigUUID)

		if err != nil {
			return nil, err
		}

		data, err := dev.GetCharByUUID(HumidityDataUUID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("Cannot find HumidityData characteristic " + HumidityDataUUID)
		}

		period, err := dev.GetCharByUUID(HumidityPeriodUUID)
		if err != nil {
			return nil, err
		}
		if period == nil {
			return nil, errors.New("Cannot find HumidityPeriod characteristic " + HumidityPeriodUUID)
		}

		return &HumiditySensor{tag, cfg, data, period}, err
	})

	return i.(*HumiditySensor), err
}

//HumiditySensor struct
type HumiditySensor struct {
	tag    *SensorTag
	cfg    *gatt.GattCharacteristic1
	data   *gatt.GattCharacteristic1
	period *gatt.GattCharacteristic1
}

//GetName return the sensor name
func (s HumiditySensor) GetName() string {
	return "humidity"
}

//Enable humidity measurements
func (s *HumiditySensor) Enable() error {
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

//Disable humidity measurements
func (s *HumiditySensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]interface{})
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

// IsEnabled check if humidity measurements are enabled
func (s *HumiditySensor) IsEnabled() (bool, error) {
	options := make(map[string]interface{})

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

//IsNotifying check if humidity sensor is notyfing
func (s *HumiditySensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.Value().(bool), nil
}

//Read value from the humidity sensor
func (s *HumiditySensor) Read() (float64, error) {
	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]interface{})
	b, err := s.data.ReadValue(options)

	if err != nil {
		return 0, err
	}

	humid := binary.LittleEndian.Uint16(b[2:])
	humidValue := calcHumidLocal(uint16(humid))

	return humidValue, err
}

//StartNotify enable DataChannel for humidity
func (s *HumiditySensor) StartNotify() error {

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

			b1 := prop.Value.([]byte)

			humid := binary.LittleEndian.Uint16(b1[2:])
			humidityValue := calcHumidLocal(uint16(humid))

			temperature := binary.LittleEndian.Uint16(b1[0:2])
			tempValue := calcTmpFromHumidSensor(uint16(temperature))
			dataEvent := SensorTagDataEvent{
				Device:            s.tag.Device1,
				SensorType:        "humidity",
				HumidityValue:     humidityValue,
				HumidityUnit:      "%RH",
				HumidityTempValue: tempValue,
				HumidityTempUnit:  "C",
				SensorID:          s.tag.Properties.Address,
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

//StopNotify disable DataChannel for humidity sensor
func (s *HumiditySensor) StopNotify() error {

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
