package override

func GetPropertyTypeOverride(iface string, property string) (string, bool) {
	if props, ok := PropertyTypes[iface]; ok {
		val, ok := props[property]
		return val, ok
	}
	return "", false
}

var PropertyTypes = map[string]map[string]string{
	"org.bluez.Device1": map[string]string{
		// "ServiceData":      "map[string]dbus.Variant",
		"ManufacturerData": "map[uint16]dbus.Variant",
	},
	"org.bluez.GattCharacteristic1": map[string]string{
		"Value": "[]byte `dbus:emit`",
	},
	"org.bluez.GattDescriptor1": map[string]string{
		"Value": "[]byte `dbus:emit`",
	},
	"org.bluez.GattService1": map[string]string{
		"Characteristics": "[]dbus.ObjectPath `dbus:emit`",
	},
	"org.bluez.LEAdvertisement1": map[string]string{
		"Data": "map[uint8][]byte",
	},
}
