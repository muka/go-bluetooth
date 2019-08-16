package api

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	log "github.com/sirupsen/logrus"
)

// NewDBusProperties create a new instance
func NewDBusProperties(conn *dbus.Conn) (*DBusProperties, error) {

	o := &DBusProperties{
		conn:        conn,
		props:       make(map[string]bluez.Properties),
		propsConfig: make(map[string]map[string]*prop.Prop),
	}

	err := o.parseProperties()
	return o, err
}

// DBus Properties interface implementation
type DBusProperties struct {
	conn        *dbus.Conn
	props       map[string]bluez.Properties
	propsConfig map[string]map[string]*prop.Prop
	instance    *prop.Properties
}

func (p *DBusProperties) parseTag(t *structs.Struct, field *structs.Field, conf *prop.Prop, tag string) bool {

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
				return true
			} else {

				checkField, ok := t.FieldOk(tagValue)
				if !ok {
					log.Warnf("%s: field not found,  is it avaialable?", tagValue)
					return false
				}
				if !checkField.IsExported() {
					log.Warnf("%s: field must be exported. (add a tag `ignore` to avoid exposing it as property)", tagValue)
					return false
				}

				varKind := checkField.Kind()
				if varKind != reflect.Bool {
					log.Warnf("%s: ignore tag expect a bool property to check, %s given", tagValue, varKind)
					return false
				}

				if checkField.Value().(bool) {
					return true
				}
			}
		}

		// check if empty
		if tagKey == "omitEmpty" {
			v := reflect.ValueOf(conf.Value)
			if !v.IsValid() || isEmptyValue(v) {
				return true
			}
		}

		switch tagKey {
		case "emit":
			conf.Emit = prop.EmitTrue
			conf.Writable = true
			break
		case "invalidates":
			conf.Emit = prop.EmitInvalidates
			conf.Writable = true
			break
		case "writable":
			conf.Writable = true
			break
		default:
			t := reflect.TypeOf(p)
			m, ok := t.MethodByName(tagKey)
			if ok {
				conf.Writable = true
				conf.Callback = m.Func.Interface().(func(*prop.Change) *dbus.Error)
			}
		}
	}

	return false
}

func (p *DBusProperties) parseProperties() error {
	for iface, ifaceVal := range p.props {

		if _, ok := p.propsConfig[iface]; !ok {
			p.propsConfig[iface] = make(map[string]*prop.Prop)
		}

		t := structs.New(ifaceVal)
		for _, field := range t.Fields() {

			if !field.IsExported() {
				continue
			}

			if _, ok := field.Value().(dbus.ObjectPath); ok && field.IsZero() {
				// log.Debugf("parseProperties: skip empty ObjectPath %s", field.Name())
				continue
			}

			propConf := &prop.Prop{
				Value:    field.Value(),
				Emit:     prop.EmitFalse,
				Writable: false,
				Callback: p.onChange,
			}

			tag := field.Tag("dbus")
			if tag != "" {
				skip := p.parseTag(t, field, propConf, tag)
				// log.Printf("%t %s", skip, tag)
				if skip {
					continue
				}
			}

			// log.Debugf("parseProperties: %s: `%s` %v", field.Name(), tag, propConf)
			p.propsConfig[iface][field.Name()] = propConf
		}
	}
	return nil
}

func (p *DBusProperties) onChange(ev *prop.Change) *dbus.Error {
	if _, ok := p.propsConfig[ev.Iface]; ok {
		if conf, ok := p.propsConfig[ev.Iface][ev.Name]; ok {
			if conf.Writable {
				log.Debugf("Set %s.%s", ev.Iface, ev.Name)
				prop := p.props[ev.Iface]
				s := structs.New(prop)
				err := s.Field(ev.Name).Set(ev.Value)
				if err != nil {
					log.Errorf("Failed to set %s.%s: %s", ev.Iface, ev.Name, err.Error())
					return &profile.ErrRejected
				}
			}
		}
	}
	return nil
}

//Instance return the props instance
func (p *DBusProperties) Instance() *prop.Properties {
	return p.instance
}

//Introspection return the props instance
func (p *DBusProperties) Introspection(iface string) []introspect.Property {
	return p.instance.Introspection(iface)
}

//Expose expose the properties interface
func (p *DBusProperties) Expose(path dbus.ObjectPath) {
	p.instance = prop.New(p.conn, path, p.propsConfig)
}

//AddProperties add a property set
func (p *DBusProperties) AddProperties(iface string, props bluez.Properties) error {
	p.props[iface] = props
	return p.parseProperties()
}

//RemoveProperties remove a property set
func (p *DBusProperties) RemoveProperties(iface string) {
	if _, ok := p.props[iface]; ok {
		delete(p.props, iface)
	}
	if _, ok := p.propsConfig[iface]; ok {
		delete(p.propsConfig, iface)
	}
}

// check for empy value, from go encoding/json
// https://github.com/golang/go/blob/master/src/encoding/json/encode.go#L318
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
