package api

import (
	"reflect"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/emitter"
	"github.com/muka/device-manager/util"
)

// NewDevice creates a new Device
func NewDevice(path string) *Device {
	d := new(Device)
	d.Path = path
	d.client = bluez.NewDevice1(path)
	return d
}

// ParseDevice parse a Device from a ObjectManager map
func ParseDevice(path dbus.ObjectPath, propsMap map[string]dbus.Variant) *Device {

	d := new(Device)
	d.Path = string(path)
	d.client = bluez.NewDevice1(d.Path)

	props := new(bluez.Device1Properties)
	util.MapToStruct(props, propsMap)
	d.client.Properties = props

	return d
}

// PropertyChanged an object to describe a changed property
type PropertyChanged struct {
	Iface      string
	Field      string
	Value      interface{}
	Properties *bluez.Device1Properties
}

func (d *Device) watchProperties() error {

	logger.Println("Registering to PropertyChanged")

	channel, err := d.client.Register()
	if err != nil {
		return err
	}

	go (func() {
		for {
			if channel == nil {
				// logger.Println("Quit goroutine")
				break
			}

			// logger.Println("Waiting for property change")
			sig := <-channel

			// logger.Println("----------------------")
			// logger.Printf("Name: %s\n", sig.Name)

			if sig.Name != bluez.PropertiesChanged {
				// logger.Printf("Skipped %s vs %s\n", sig.Name, bluez.PropertiesInterface)
				continue
			}

			for i := 0; i < len(sig.Body); i++ {
				logger.Println(reflect.TypeOf(sig.Body[i]))
				logger.Println(sig.Body[i])
			}

			// logger.Println("----------------------")

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)
			for field, val := range changes {

				// updates [*]Properties struct
				props := d.GetProperties()
				s := reflect.ValueOf(props).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						f.Set(x)
						// logger.Printf("Set props value: %s = %s\n", field, x.Interface())
					}
				}

				propChanged := PropertyChanged{string(iface), field, val.Value(), props}
				d.Emit("change", propChanged)
			}
		}
	})()

	return nil
}

func (d *Device) unwatchProperties() error {
	return d.client.Unregister()
}

//Device return an API to interact with a DBus device
type Device struct {
	Path string
	// Properties *bluez.Device1Properties
	client *bluez.Device1
}

//GetClient return a DBus Device1 interface client
func (d *Device) GetClient() *bluez.Device1 {
	return d.client
}

//GetProperties return the properties for the device
func (d *Device) GetProperties() *bluez.Device1Properties {
	if d.client.Properties == nil {
		d.client.GetProperties()
	}
	return d.client.Properties
}

//On register callback for event
func (d *Device) On(name string, fn emitter.Callback) {
	switch name {
	case "change":
		d.watchProperties()
		break
	}
	emitter.On(d.Path+"."+name, fn)
}

//Off unregister callback for event
func (d *Device) Off(name string) {
	switch name {
	case "change":
		d.unwatchProperties()
		break
	}
	emitter.Off(d.Path + "." + name)
}

//Emit an event
func (d *Device) Emit(name string, data interface{}) {
	emitter.Emit(d.Path+"."+name, data)
}
