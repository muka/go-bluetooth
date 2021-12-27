package override

var typesMap = map[string]string{
	//mesh-api object node, array{byte, array{(uint16, dict)}} configuration Attach(object app_root, uint64 token)
	"object node, array{byte, array{(uint16, dict)}} configuration": "dbus.ObjectPath, []ConfigurationItem",

	"array{(uint16 id, dict caps)}": "[]ConfigurationItem",
	// mesh-api array{(uint16, uint16)} VendorModels [read-only]
	"array{(uint16, uint16)}": "[]VendorItem",

	"array{(uint16 vendor, uint16 id, dict options)}": "[]VendorOptionsItem",
	// obex-api array{string vcard, string name} List(dict filters)
	"array{string vcard, string name}": "[]VCardItem",
	// obex-api
	"object, dict": "dbus.ObjectPath, map[string]interface{}",
	// obex-api array{object, dict} ListMessages(string folder, dict filter)
	"array{object, dict}": "[]Message",
	// gatt AcquireWrite
	"fd, uint16": "dbus.UnixFD, uint16",
	// media-api array{objects, properties} ListItems(dict filter)
	"array{objects, properties}": "[]Item",
	// media-api fd, uint16, uint16 Acquire()
	"fd, uint16, uint16": "dbus.UnixFD, uint16, uint16",
	// advertisement_monitor array{(uint8, uint8, array{byte})} Patterns [read-only, optional]
	"array{(uint8, uint8, array{byte})}": "[]Pattern",
	// advertisement_monitor Uint/Int with uppercase
	"Uint16": "uint16",
	"Int16":  "int16",
}

//MapType map a raw type literal otherwise difficult to parse
func MapType(rawtype string) (string, bool) {
	res, ok := typesMap[rawtype]
	return res, ok
}
