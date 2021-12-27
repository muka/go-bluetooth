package override

func GetPropertiesOverride(iface string) (map[string]string, bool) {
	if props, ok := PropertyTypes[iface]; ok {
		return props, ok
	}
	return map[string]string{}, false
}

var PropertyTypes = map[string]map[string]string{
	"org.bluez.Device1": {
		"ServiceData":      "map[string]interface{}",
		"ManufacturerData": "map[uint16]interface{}",
	},
	"org.bluez.GattCharacteristic1": {
		"Value":          "[]byte `dbus:\"emit\"`",
		"Descriptors":    "[]dbus.ObjectPath",
		"WriteAcquired":  "bool `dbus:\"ignore\"`",
		"NotifyAcquired": "bool `dbus:\"ignore\"`",
	},
	"org.bluez.GattDescriptor1": {
		"Value":          "[]byte `dbus:\"emit\"`",
		"Characteristic": "dbus.ObjectPath",
	},
	"org.bluez.GattService1": {
		"Characteristics": "[]dbus.ObjectPath `dbus:\"emit\"`",
		"Includes":        "[]dbus.ObjectPath `dbus:\"omitEmpty\"`",
		"Device":          "dbus.ObjectPath `dbus:\"ignore=IsService\"`",
		"IsService":       "bool `dbus:\"ignore\"`",
	},
	"org.bluez.LEAdvertisement1": {
		// dbus type: (yv) dict of byte variant (array of bytes)
		"Data": "map[byte]interface{}",
		// dbus type: (qv) dict of uint16 variant (array of bytes)
		"ManufacturerData": "map[uint16]interface{}",
		// dbus type: (s[v]) dict of string variant (array of bytes)
		"ServiceData": "map[string]interface{}",
		// SecondaryChannel, if set on 5.54 cause a parsing exception
		"SecondaryChannel": "string `dbus:\"omitEmpty\"`",
	},
	"org.bluez.AdvertisementMonitor1": {
		// array{(uint8, uint8, array{byte})}
		"Patterns": "[]Pattern",
	},
	"org.bluez.MediaPlayer1": {
		"Track": "Track",
	},
}
