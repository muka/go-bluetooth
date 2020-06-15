package api

import (
	"github.com/fatih/structs"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/props"
	log "github.com/sirupsen/logrus"
)

// NewDBusProperties create a new instance
func NewDBusProperties(conn *dbus.Conn) (*DBusProperties, error) {

	o := &DBusProperties{
		conn:        conn,
		props:       make(map[string]bluez.Properties),
		propsConfig: make(map[string]map[string]*props.PropInfo),
	}

	err := o.parseProperties()
	return o, err
}

// DBus Properties interface implementation
type DBusProperties struct {
	conn        *dbus.Conn
	props       map[string]bluez.Properties
	propsConfig map[string]map[string]*props.PropInfo
	instance    *prop.Properties
}

func (p *DBusProperties) parseProperties() error {
	for iface, ifaceVal := range p.props {
		if _, ok := p.propsConfig[iface]; !ok {
			p.propsConfig[iface] = make(map[string]*props.PropInfo)
		}
		p.propsConfig[iface] = props.ParseProperties(ifaceVal)
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
	res := p.instance.Introspection(iface)
	// log.Debug("Introspect", res)
	return res
}

//Expose expose the properties interface
func (p *DBusProperties) Expose(path dbus.ObjectPath) {
	propsConfig := make(map[string]map[string]*prop.Prop)
	for iface1, props1 := range p.propsConfig {
		propsConfig[iface1] = make(map[string]*prop.Prop)
		for k, v := range props1 {
			if v.Skip {
				continue
			}
			propsConfig[iface1][k] = &v.Prop
		}
	}

	p.instance = prop.New(p.conn, path, propsConfig)
}

//AddProperties add a property set
func (p *DBusProperties) AddProperties(iface string, props bluez.Properties) error {
	p.props[iface] = props
	return p.parseProperties()
}

//RemoveProperties remove a property set
func (p *DBusProperties) RemoveProperties(iface string) {
	delete(p.props, iface)
	delete(p.propsConfig, iface)
}
