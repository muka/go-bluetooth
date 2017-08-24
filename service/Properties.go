package service

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewProperties create a new instance
func NewProperties(conn *dbus.Conn) (*Properties, error) {

	o := &Properties{
		conn:  conn,
		props: make(map[string]bluez.Properties),
	}

	return o, nil
}

// Properties interface implementation
type Properties struct {
	conn  *dbus.Conn
	props map[string]bluez.Properties
}

//AddProperties add a property set
func (p *Properties) AddProperties(iface string, props bluez.Properties) {
	p.props[iface] = props
}

//RemoveProperties remove a property set
func (p *Properties) RemoveProperties(iface string) {
	if _, ok := p.props[iface]; ok {
		delete(p.props, iface)
	}
}

//Get return a property by name
func (p *Properties) Get(iface string, propertyName string) (interface{}, error) {
	if prop, ok := p.props[iface]; ok {
		list := prop.ToMap()
		if val, ok := list[propertyName]; ok {
			return val, nil
		}
	}
	return nil, errors.New("Property not found")
}

//Set the value of a property by name
func (p *Properties) Set(iface string, propertyName string, val dbus.Variant) error {
	if prop, ok := p.props[iface]; ok {
		list := prop.ToMap()
		if _, ok := list[propertyName]; ok {
			s := structs.New(prop)
			s.Field(propertyName).Set(val.Value())
			return nil
		}
	}
	return errors.New("Property not found")
}

//GetAll return all the properties by interface name
func (p *Properties) GetAll(iface string) (map[string]interface{}, error) {
	if prop, ok := p.props[iface]; ok {
		list := prop.ToMap()
		return list, nil
	}
	return nil, errors.New("Interface not found")
}
