package sensortag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

//.....getting config,data,period characteristics for TEMPERATURE sensor............
func newTemperatureSensor(tag *SensorTag) (*TemperatureSensor, error) {

	dev := tag.Device

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

	return i.(*TemperatureSensor), err
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

//Read value from the sensor
func (s *TemperatureSensor) Read() (float64, error) {

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

//StartNotify enable temperature DataChannel
func (s *TemperatureSensor) StartNotify() error {

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
					SensorID:         s.tag.Properties.Address,
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
