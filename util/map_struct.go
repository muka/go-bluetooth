package util

import (
	"fmt"
	"reflect"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

func AssignMapVariantToInterface(mapVal reflect.Value, mapVariant reflect.Value) (bool, error) {

	if mapVal.Type().Kind() != reflect.Map {
		return false, nil
	}
	if mapVariant.Type().Kind() != reflect.Map {
		return false, nil
	}

	// map[*]interface{}
	mapValT := mapVal.Type()
	mapValKey := mapValT.Key()
	mapValValue := mapValT.Elem()

	// map[*]dbus.Variant
	mapVariantT := mapVariant.Type()
	mapVariantKey := mapVariantT.Key()
	mapVariantValue := mapVariantT.Elem()

	// keys must match
	if mapValKey.Kind() != mapVariantKey.Kind() {
		return false, fmt.Errorf(
			"Cannot set values on different types: map[%s] with map[%s]",
			mapValKey.Kind().String(),
			mapVariantKey.Kind().String(),
		)
	}

	// receiving value is interface{}
	if mapValValue.Kind() != reflect.Interface {
		log.Debugf("val is not interface")
		return false, nil
	}

	// source value is dbus.Variant
	if mapVariantValue.Kind() != reflect.TypeOf(dbus.Variant{}).Kind() {
		log.Debugf("mapVariant value is not variant")
		return false, nil
	}

	mapValInstanceType := reflect.MapOf(mapValKey, mapValValue)
	mapValInstance := reflect.MakeMapWithSize(mapValInstanceType, mapVariant.Len())

	for _, key := range mapVariant.MapKeys() {
		variantInnerValue := mapVariant.MapIndex(key).MethodByName("Value").Call([]reflect.Value{})
		// log.Debugf("variantInnerValue %++v", variantInnerValue)
		mapValInstance.SetMapIndex(key, variantInnerValue[0])
	}

	mapVal.Set(mapValInstance)

	return true, nil
}

func mapStructField(obj interface{}, name string, value dbus.Variant) error {

	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		log.Warnf("MapToStruct: invalid field detected %T.%s", obj, name)
		return nil
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set value for: %s", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value.Value())

	if structFieldType == val.Type() {
		structFieldValue.Set(val)
		return nil
	}

	// log.Debugf("structFieldType %++v", structFieldType)
	// log.Debugf("val.Type() %++v", val.Type())

	if val.Type().Kind() == reflect.Map {

		if structFieldType.Kind() == reflect.Struct ||
			(structFieldType.Kind() == reflect.Pointer && structFieldType.Elem().Kind() == reflect.Struct) {

			variantMap, ok := value.Value().(map[string]dbus.Variant)
			if ok {
				if structFieldType.Kind() == reflect.Pointer {
					if structFieldValue.IsZero() {
						return fmt.Errorf("Pointer field %s: uninitialized", name)
					}
					return MapToStruct(structFieldValue.Interface(), variantMap)
				}
				return MapToStruct(structFieldValue.Addr().Interface(), variantMap)
			}
		}

		structVal := structFieldType.Elem()
		structKey := structFieldType.Key()

		// mapVal := val.Type().Elem()
		mapKey := val.Type().Key()

		if mapKey.Kind() != structKey.Kind() {
			return fmt.Errorf("Field %s: map key mismatching values object=%s props=%s", name, structKey.Kind(), mapKey.Kind())
		}

		// Assign value if signature is map[*]interface{}
		if structVal.Kind() == reflect.Interface {

			val1MapType := reflect.MapOf(structKey, structVal)
			val1Map := reflect.MakeMapWithSize(val1MapType, val.Len())

			for _, key := range val.MapKeys() {
				val1Map.SetMapIndex(key, val.MapIndex(key))
			}

			structFieldValue.Set(val1Map)
			return nil
		}

	}

	if val.Type().Kind() == reflect.Array {
		log.Warn("@TODO type array to interface{} is not implemented")
	}

	return fmt.Errorf("Mismatching types for field=%s object=%s props=%s", name, structFieldType, val.Type())
}

// MapToStruct converts a map[string]dbus.Variant to a struct
func MapToStruct(s interface{}, m map[string]dbus.Variant) error {
	for k, v := range m {
		err := mapStructField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// StructToMap converts a struct to a map[string]interface{}
func StructToMap(s interface{}, m map[string]interface{}) {
	structValue := reflect.ValueOf(s).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		m[structValue.Type().Field(i).Name] = structValue.Field(i).Interface()
	}
}
