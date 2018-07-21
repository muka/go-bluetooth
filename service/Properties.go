package service

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
)

// NewProperties create a new instance
func NewProperties(conn *dbus.Conn) (*Properties, error) {

	o := &Properties{
		conn:        conn,
		props:       make(map[string]bluez.Properties),
		propsConfig: make(map[string]map[string]*prop.Prop),
	}

	err := o.parseProperties()
	return o, err
}

// Properties interface implementation
type Properties struct {
	conn        *dbus.Conn
	props       map[string]bluez.Properties
	propsConfig map[string]map[string]*prop.Prop
	instance    *prop.Properties
}

func (p *Properties) parseTag(conf *prop.Prop, tag string) {
	parts := strings.Split(tag, ",")
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
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
			m, ok := t.MethodByName(parts[i])
			if ok {
				conf.Writable = true
				conf.Callback = m.Func.Interface().(func(*prop.Change) *dbus.Error)
			}
		}
	}
}

func (p *Properties) parseProperties() error {
	for iface, ifaceVal := range p.props {

		if _, ok := p.propsConfig[iface]; !ok {
			p.propsConfig[iface] = make(map[string]*prop.Prop)
		}

		t := structs.New(ifaceVal)
		for _, field := range t.Fields() {

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
				p.parseTag(propConf, tag)
			}

			// log.Debugf("parseProperties: %s: `%s` %v", field.Name(), tag, propConf)
			p.propsConfig[iface][field.Name()] = propConf
		}
	}
	return nil
}

func (p *Properties) onChange(ev *prop.Change) *dbus.Error {
	if _, ok := p.propsConfig[ev.Iface]; ok {
		if conf, ok := p.propsConfig[ev.Iface][ev.Name]; ok {
			if conf.Writable {
				log.Debugf("Set %s.%s", ev.Iface, ev.Name)
				prop := p.props[ev.Iface]
				s := structs.New(prop)
				err := s.Field(ev.Name).Set(ev.Value)
				if err != nil {
					log.Errorf("Failed to set %s.%s: %s", ev.Iface, ev.Name, err.Error())
					return DbusErr
				}
			}
		}
	}
	return nil
}

//Instance return the props instance
func (p *Properties) Instance() *prop.Properties {
	return p.instance
}

//Introspection return the props instance
func (p *Properties) Introspection(iface string) []introspect.Property {
	return p.instance.Introspection(iface)
}

//Expose expose the properties interface
func (p *Properties) Expose(path dbus.ObjectPath) {
	p.instance = prop.New(p.conn, path, p.propsConfig)
}

//AddProperties add a property set
func (p *Properties) AddProperties(iface string, props bluez.Properties) error {
	p.props[iface] = props
	return p.parseProperties()
}

//RemoveProperties remove a property set
func (p *Properties) RemoveProperties(iface string) {
	if _, ok := p.props[iface]; ok {
		delete(p.props, iface)
	}
	if _, ok := p.propsConfig[iface]; ok {
		delete(p.propsConfig, iface)
	}
}
