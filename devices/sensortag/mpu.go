package sensortag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
)

//getting config,data,period characteristics for Humidity sensor
func newMpuSensor(tag *SensorTag) (*MpuSensor, error) {

	dev := tag.Device

	//accelerometer,magnetometer,gyroscope
	MpuConfigUUID, err := getUUID("MPU9250_CONFIG_UUID")
	if err != nil {
		return nil, err
	}
	MpuDataUUID, err := getUUID("MPU9250_DATA_UUID")
	if err != nil {
		return nil, err
	}
	MpuPeriodUUID, err := getUUID("MPU9250_PERIOD_UUID")
	if err != nil {
		return nil, err
	}

	i, err := retryCall(DefaultRetry, DefaultRetryWait, func() (interface{}, error) {

		cfg, err := dev.GetCharByUUID(MpuConfigUUID)
		if err != nil {
			return nil, err
		}

		data, err := dev.GetCharByUUID(MpuDataUUID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("Cannot find MpuData characteristic " + MpuDataUUID)
		}

		period, err := dev.GetCharByUUID(MpuPeriodUUID)
		if err != nil {
			return nil, err
		}
		if period == nil {
			return nil, errors.New("Cannot find MpuPeriod characteristic " + MpuPeriodUUID)
		}

		return &MpuSensor{tag, cfg, data, period}, err
	})

	return i.(*MpuSensor), err
}

//MpuSensor structure
type MpuSensor struct {
	tag    *SensorTag
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

//GetName return's the sensor name
func (s MpuSensor) GetName() string {
	return "Ac-Mg-Gy"
}

//Enable mpuSensors measurements
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

//Disable mpuSensors measurements
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

//IsEnabled check if mpu measurements are enabled
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

//IsNotifying check if mpu sensors are notifying
func (s *MpuSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.(bool), nil
}

//Read value from the mpu sensors
func (s *MpuSensor) Read() (float64, error) {

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

//StartNotify enable mpuDataChannel
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

				//... calculate Gyroscope...........

				mpuXg := binary.LittleEndian.Uint16(b1[0:2])
				mpuYg := binary.LittleEndian.Uint16(b1[2:4])
				mpuZg := binary.LittleEndian.Uint16(b1[4:6])

				mpuGyX, mpuGyY, mpuGyZ := calcMpuGyroscope(uint16(mpuXg), uint16(mpuYg), uint16(mpuZg))
				mpuGyroscope = fmt.Sprint(mpuGyX, " , ", mpuGyY, " , ", mpuGyZ)

				//... calculate Accelerometer.......

				mpuXa := binary.LittleEndian.Uint16(b1[6:8])
				mpuYa := binary.LittleEndian.Uint16(b1[8:10])
				mpuZa := binary.LittleEndian.Uint16(b1[10:12])

				mpuAcX, mpuAcY, mpuAcZ := calcMpuAccelerometer(uint16(mpuXa), uint16(mpuYa), uint16(mpuZa))
				mpuAccelerometer = fmt.Sprint(mpuAcX, " , ", mpuAcY, " , ", mpuAcZ)

				//... calculate Magnetometer.......

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
					SensorID:              macAddress,
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

//StopNotify disable DataChannel for mpu sensors
func (s *MpuSensor) StopNotify() error {

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
