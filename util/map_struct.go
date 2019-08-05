package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/godbus/dbus"
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
		return fmt.Errorf("Provided value type didn't match obj field type. field=%s expected=%s actual=%s", name, structFieldType, val.Type())
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

// StructToMap converts a struct to a map[string]interface{}
func StructToMap(s interface{}, m map[string]interface{}) error {

	structValue := reflect.ValueOf(s).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		m[structValue.Type().Field(i).Name] = structValue.Field(i).Interface()
	}

	return nil
}
