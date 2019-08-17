package props

import (
	"github.com/muka/go-bluetooth/bluez"
)

// Convert a struct to map applying options from struct tag
func ToMap(a bluez.Properties) map[string]interface{} {

	propsInfo := ParseProperties(a)

	res := make(map[string]interface{})
	for name, info := range propsInfo {
		if info.Skip {
			continue
		}
		res[name] = info.Value
	}

	return res
}
