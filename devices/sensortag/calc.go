package sensortag

import "math"

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor
var calcBarometricPressure = func(raw uint32) float64 {
	//barometric pressure......
	pressureMask := (int(raw) >> 8) & 0x00ffffff
	return float64(pressureMask) / 100.0
}

var calcBarometricTemperature = func(raw uint32) float64 {
	//TEMPERATURE calibiration data coming from barometric sensor
	tempMask := int(raw) & 0x00ffffff
	return float64(tempMask) / 100.0
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor
var calcHumidLocal = func(raw uint16) float64 {
	//humidity calibiration.........
	return float64(raw) * 100 / 65536.0
}

var calcTmpFromHumidSensor = func(raw uint16) float64 {
	//TEMPERATURE calibiration for data coming from humidity sensor
	return -40 + ((165 * float64(raw)) / 65536.0)
}

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor
var calcTmpLocal = func(raw uint16) float64 {
	//ambient temperature calberation
	return float64(raw) / 128.0
}

/* Conversion algorithm for target temperature */
var calcTmpTarget = func(raw uint16) float64 {
	//..object temperature caliberation...........
	return float64(raw) / 128.0
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

// Port from http://processors.wiki.ti.com/index.php/SensorTag_User_Guide#IR_Temperature_Sensor

var calcLuxometer = func(raw uint16) float64 {

	exponent := (int(raw) & 0xF000) >> 12
	mantissa := (int(raw) & 0x0FFF)
	exp := float64(exponent)
	man := float64(mantissa)
	flLux := man * math.Pow(2, exp) / 100
	return float64(flLux)
}
