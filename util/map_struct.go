package util

import (
	"errors"
	"github.com/godbus/dbus"
	"reflect"
)

func mapStructField(obj interface{}, name string, value dbus.Variant) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New("No such field: " + name + " in obj")
	}

	if !structFieldValue.CanSet() {
		return errors.New("Cannot set " + name + " field value")
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value.Value())
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

// MapToStruct converts a map[string]interface{} to a struct
func MapToStruct(s interface{}, m map[string]dbus.Variant) error {
	for k, v := range m {
		err := mapStructField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
