package override

func GetPropertiesOverride(iface string) (map[string]string, bool) {
	if props, ok := PropertyTypes[iface]; ok {
		return props, ok
	}
	return map[string]string{}, false
}

var PropertyTypes = map[string]map[string]string{
	"org.bluez.Device1": map[string]string{
		"ServiceData":      "map[string]interface{}",
		"ManufacturerData": "map[uint16]interface{}",
	},
	"org.bluez.GattCharacteristic1": map[string]string{
		"Value":          "[]byte `dbus:\"emit\"`",
		"Descriptors":    "[]dbus.ObjectPath",
		"WriteAcquired":  "bool `dbus:\"ignore\"`",
		"NotifyAcquired": "bool `dbus:\"ignore\"`",
	},
	"org.bluez.GattDescriptor1": map[string]string{
		"Value":          "[]byte `dbus:\"emit\"`",
		"Characteristic": "dbus.ObjectPath",
	},
	"org.bluez.GattService1": map[string]string{
		"Characteristics": "[]dbus.ObjectPath `dbus:\"emit\"`",
		"Includes":        "[]dbus.ObjectPath `dbus:\"omitEmpty\"`",
		"Device":          "dbus.ObjectPath `dbus:\"ignore=IsService\"`",
		"IsService":       "bool `dbus:\"ignore\"`",
	},
	"org.bluez.LEAdvertisement1": map[string]string{
		// dbus type: (yv) dict of byte variant (array of bytes)
		"Data": "map[byte]interface{}",
		// dbus type: (qv) dict of uint16 variant (array of bytes)
		"ManufacturerData": "map[uint16]interface{}",
		// dbus type: (s[v]) dict of string variant (array of bytes)
		"ServiceData": "map[string]interface{}",
		// SecondaryChannel, if set on 5.54 cause a parsing exception
		"SecondaryChannel": "string `dbus:\"omitEmpty\"`",
	},
}
