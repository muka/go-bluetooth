package override

var typesMap = map[string]string{
	// mesh Network1.Attach()
	// object node, array{byte, array{(uint16, dict)}} configuration
	"object node, array{byte, array{(uint16, dict)}} configuration": "dbus.ObjectPath, []map[byte][]map[uint16]map[string]interface{}",
	// mesh Element1
	// array{(uint16, uint16)}
	"array{(uint16, uint16)}": "[]map[uint16]uint16",
	// obex org.bluez.obex.PhonebookAccess1.List
	"array{string vcard, string name}": "[]map[string]string",
	// obex
	"object, dict": "dbus.ObjectPath, map[string]interface{}",
	// obext List
	"array{object, dict}": "[]map[dbus.ObjectPath]map[string]interface{}",
	// gatt AcquireWrite
	"fd, uint16": "dbus.UnixFD, uint16",
}

//MapType map a raw type literal otherwise difficult to parse
func MapType(rawtype string) (string, bool) {
	res, ok := typesMap[rawtype]
	return res, ok
}
