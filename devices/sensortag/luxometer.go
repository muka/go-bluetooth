package sensortag

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//Luxometer Sensor..

//getting config,data,period characteristics for luxometer sensor

func newLuxometerSensor(tag *SensorTag) (*LuxometerSensor, error) {

	dev := tag.Device1

	LuxometerConfigUUID, err := getUUID("LUXOMETER_CONFIG_UUID")
	if err != nil {
		return nil, err
	}
	LuxometerDataUUID, err := getUUID("LUXOMETER_DATA_UUID")
	if err != nil {
		return nil, err
	}
	LuxometerPeriodUUID, err := getUUID("LUXOMETER_PERIOD_UUID")
	if err != nil {
		return nil, err
	}

	i, err := retryCall(DefaultRetry, DefaultRetryWait, func() (interface{}, error) {

		cfg, err := dev.GetCharByUUID(LuxometerConfigUUID)
		if err != nil {
			return nil, err
		}

		data, err := dev.GetCharByUUID(LuxometerDataUUID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("Cannot find LuxometerDataUUID  characteristic " + LuxometerDataUUID)
		}

		period, err := dev.GetCharByUUID(LuxometerPeriodUUID)
		if err != nil {
			return nil, err
		}
		if period == nil {
			return nil, errors.New("Cannot find LuxometerPeriodUUID  characteristic " + LuxometerPeriodUUID)
		}

		return &LuxometerSensor{tag, cfg, data, period}, err
	})

	return i.(*LuxometerSensor), err
}

//LuxometerSensor sensor structure
type LuxometerSensor struct {
	tag    *SensorTag
	cfg    *gatt.GattCharacteristic1
	data   *gatt.GattCharacteristic1
	period *gatt.GattCharacteristic1
}

//GetName return the sensor name
func (s LuxometerSensor) GetName() string {
	return "luxometer"
}

//Enable LuxometerSensor measurements
func (s *LuxometerSensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]interface{})
	err = s.cfg.WriteValue([]byte{1}, options)
	if err != nil {
		return err
	}
	return nil
}

//Disable LuxometerSensor measurements
func (s *LuxometerSensor) Disable() error {
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

//IsEnabled check if LuxometerSensor measurements are enabled
func (s *LuxometerSensor) IsEnabled() (bool, error) {
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

//IsNotifying check if LuxometerSensor sensors are Notifying
func (s *LuxometerSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.Value().(bool), nil
}

//Read value from the LuxometerSensor sensors
func (s *LuxometerSensor) Read() (float64, error) {

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

//StartNotify enable LuxometerSensorDataChannel
func (s *LuxometerSensor) StartNotify(macAddress string) error {

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
			luxometer := binary.LittleEndian.Uint16(b1[0:])
			luxometerValue := calcLuxometer(uint16(luxometer))

			dataEvent := SensorTagDataEvent{
				Device:         s.tag.Device1,
				SensorType:     "luxometer",
				LuxometerValue: luxometerValue,
				LuxometerUnit:  "candela",
				SensorID:       macAddress,
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

//StopNotify disable Luxometer Sensor DataChannel
func (s *LuxometerSensor) StopNotify() error {

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
