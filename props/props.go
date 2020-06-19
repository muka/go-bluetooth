package props

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

type PropInfo struct {
	prop.Prop
	Skip bool
}

func ParseProperties(propertyVal bluez.Properties) map[string]*PropInfo {

	t := structs.New(propertyVal)

	res := map[string]*PropInfo{}

	for _, field := range t.Fields() {

		if !field.IsExported() {
			continue
		}

		// if _, ok := field.Value().(dbus.ObjectPath); ok && field.IsZero() {
		// 	log.Debugf("parseProperties: skip empty ObjectPath %s", field.Name())
		// 	continue
		// }

		propInfo := new(PropInfo)
		propInfo.Value = field.Value()

		res[field.Name()] = propInfo

		tag := field.Tag("dbus")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		for i := 0; i < len(parts); i++ {

			tagKey := parts[i]
			tagValue := ""
			if strings.Contains(parts[i], "=") {
				subpts := strings.Split(parts[i], "=")
				tagKey = subpts[0]
				tagValue = strings.Join(subpts[1:], "=")
			}

			if tagKey == "ignore" {
				if tagValue == "" {
					propInfo.Skip = true
				} else {

					checkField, ok := t.FieldOk(tagValue)
					if !ok {
						log.Warnf("%s: field not found,  is it avaialable?", tagValue)
						continue
					}
					if !checkField.IsExported() {
						log.Warnf("%s: field must be exported. (add a tag `ignore` to avoid exposing it as property)", tagValue)
						continue
					}

					varKind := checkField.Kind()
					if varKind != reflect.Bool {
						log.Warnf("%s: ignore tag expect a bool property to check, %s given", tagValue, varKind)
						continue
					}

					if checkField.Value().(bool) {
						propInfo.Skip = true
					}
				}
			}

			// check if empty
			if tagKey == "omitEmpty" {
				if field.IsZero() {
					propInfo.Skip = true
				}
			}

			switch tagKey {
			case "emit":
				propInfo.Emit = prop.EmitTrue
				propInfo.Writable = true
			case "invalidates":
				propInfo.Emit = prop.EmitInvalidates
				propInfo.Writable = true
			case "writable":
				propInfo.Writable = true
			default:
				t := reflect.TypeOf(propertyVal)
				m, ok := t.MethodByName(tagKey)
				if ok {
					propInfo.Writable = true
					propInfo.Callback = m.Func.Interface().(func(*prop.Change) *dbus.Error)
				}
			}
		}

	}

	return res
}
