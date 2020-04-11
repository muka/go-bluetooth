package override

type ConstructorOverride struct {
	AdapterAsArgument bool
}

var constructorOverrides = map[string][]ConstructorOverride{
	"org.bluez.Adapter1": {
		ConstructorOverride{
			AdapterAsArgument: true,
		},
	},
	"org.bluez.GattManager1": {
		ConstructorOverride{
			AdapterAsArgument: true,
		},
	},
	"org.bluez.LEAdvertisingManager1": {
		ConstructorOverride{
			AdapterAsArgument: true,
		},
	},
	"org.bluez.MediaControl1": {
		ConstructorOverride{
			AdapterAsArgument: true,
		},
	},
}

func GetConstructorsOverrides(iface string) ([]ConstructorOverride, bool) {
	if val, ok := constructorOverrides[iface]; ok {
		return val, ok
	}
	return []ConstructorOverride{}, false
}
