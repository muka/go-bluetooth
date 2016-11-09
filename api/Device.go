package api

import (
	"errors"
	"reflect"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/bluez-client/bluez/profile"
	"github.com/muka/bluez-client/emitter"
	"github.com/muka/bluez-client/util"
)

// NewDevice creates a new Device
func NewDevice(path string) *Device {
	d := new(Device)
	d.Path = path
	d.client = profile.NewDevice1(path)
	return d
}

// ParseDevice parse a Device from a ObjectManager map
func ParseDevice(path dbus.ObjectPath, propsMap map[string]dbus.Variant) (*Device, error) {

	d := new(Device)
	d.Path = string(path)
	d.client = profile.NewDevice1(d.Path)

	props := new(profile.Device1Properties)
	util.MapToStruct(props, propsMap)
	c, err := d.GetClient()
	if err != nil {
		return nil, err
	}
	c.Properties = props

	return d, nil
}

func (d *Device) watchProperties() error {

	// logger.Debug("Registering to PropertyChanged")

	channel, err := d.client.Register()
	if err != nil {
		return err
	}

	go (func() {
		for {
			if channel == nil {
				// logger.Debug("Quit goroutine")
				break
			}

			// logger.Debug("Waiting for property change")
			sig := <-channel

			// logger.Debug("----------------------")
			// logger.Debug("Name: %s\n", sig.Name)

			if sig == nil {
				return
			}

			if sig.Name != bluez.PropertiesChanged {
				// logger.Debug("Skipped %s vs %s\n", sig.Name, bluez.PropertiesInterface)
				continue
			}

			// logger.Debug("Device property changed")
			// for i := 0; i < len(sig.Body); i++ {
			// 	logger.Debug(reflect.TypeOf(sig.Body[i]))
			// 	logger.Debug(sig.Body[i])
			// }

			// logger.Debug("----------------------")

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)
			for field, val := range changes {

				// updates [*]Properties struct
				props, err := d.GetProperties()

				if err != nil {
					logger.Criticalf("Exception getting properties: %v\n", err)
					return
				}

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
						// logger.Debug("Set props value: %s = %s\n", field, x.Interface())
					}
				}

				propChanged := PropertyChangedEvent{string(iface), field, val.Value(), props}
				d.Emit("changed", propChanged)
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
	// Properties *profile.Device1Properties
	client *profile.Device1
}

//GetClient return a DBus Device1 interface client
func (d *Device) GetClient() (*profile.Device1, error) {
	if d.client == nil {
		return nil, errors.New("Client not available")
	}
	return d.client, nil
}

//GetProperties return the properties for the device
func (d *Device) GetProperties() (*profile.Device1Properties, error) {
	if d == nil {
		return nil, errors.New("Empty device pointer")
	}
	c, err := d.GetClient()
	if err != nil {
		return nil, err
	}
	c.GetProperties()
	return c.Properties, nil
}

//GetProperty return a property value
func (d *Device) GetProperty(name string) (data interface{}, err error) {
	c, err := d.GetClient()
	if err != nil {
		return nil, err
	}
	val, err := c.GetProperty(name)
	if err != nil {
		return nil, err
	}
	return val.Value(), nil
}

//On register callback for event
func (d *Device) On(name string, fn Callback) {
	switch name {
	case "changed":
		d.watchProperties()
		break
	}
	logger.Debug("Listen on %s\n", d.Path+"."+name)
	emitter.On(d.Path+"."+name, func(ev emitter.Event) {
		fn(ev)
	})
}

//Off unregister callback for event
func (d *Device) Off(name string) {
	switch name {
	case "changed":
		d.unwatchProperties()
		break
	}
	emitter.Off(d.Path + "." + name)
}

//Emit an event
func (d *Device) Emit(name string, data interface{}) {
	emitter.Emit(d.Path+"."+name, data)
}

//GetService return a GattService
func (d *Device) GetService(path string) *profile.GattService1 {
	return profile.NewGattService1(path)
}

//GetChar return a GattService
func (d *Device) GetChar(path string) *profile.GattCharacteristic1 {
	return profile.NewGattCharacteristic1(path)
}
