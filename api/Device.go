package api

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/util"
)

var deviceRegistry = make(map[string]*Device)

// NewDevice creates a new Device
func NewDevice(path string) *Device {

	if _, ok := deviceRegistry[path]; ok {
		return deviceRegistry[path]
	}

	d := new(Device)
	d.Path = path
	d.client = profile.NewDevice1(path)

	d.client.GetProperties()
	d.Properties = d.client.Properties
	d.chars = make(map[dbus.ObjectPath]*profile.GattCharacteristic1, 0)

	deviceRegistry[path] = d

	// d.watchProperties()

	return d
}

//ClearDevice free a device struct
func ClearDevice(d *Device) error {

	err := d.Disconnect()
	if err != nil {
		return err
	}

	err = d.unwatchProperties()
	if err != nil {
		return err
	}

	c, err := d.GetClient()
	if err != nil {
		return err
	}

	c.Close()

	if _, ok := deviceRegistry[d.Path]; ok {
		delete(deviceRegistry, d.Path)
	}

	return nil
}

//ClearDevices clear cached devices
func ClearDevices() error {
	devices, err := GetDevices()
	if err != nil {
		return err
	}
	for _, dev := range devices {
		err := ClearDevice(&dev)
		if err != nil {
			return err
		}
	}
	return nil
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
	d.chars = make(map[dbus.ObjectPath]*profile.GattCharacteristic1, 0)

	return d, nil
}

func (d *Device) watchProperties() error {
	if d.watchPropertiesChannel != nil {
		d.unwatchProperties()
	}

	channel, err := d.client.Register()
	if err != nil {
		return err
	}
	d.watchPropertiesChannel = channel

	go (func() {
		for {

			if channel == nil {
				break
			}

			sig := <-channel

			if sig == nil {
				return
			}

			if sig.Name != bluez.PropertiesChanged {
				continue
			}
			if (fmt.Sprint(sig.Path) != d.Path) {
			    continue
			}

			// for i := 0; i < len(sig.Body); i++ {
			// log.Printf("%s -> %s\n", reflect.TypeOf(sig.Body[i]), sig.Body[i])
			// }

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)
			for field, val := range changes {

				// updates [*]Properties struct
				props := d.Properties

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
					}
				}

				propChanged := PropertyChangedEvent{string(iface), field, val.Value(), props, d}
				d.Emit("changed", propChanged)
			}
		}
	})()

	return nil
}

//Device return an API to interact with a DBus device
type Device struct {
	Path                   string
	Properties             *profile.Device1Properties
	client                 *profile.Device1
	chars                  map[dbus.ObjectPath]*profile.GattCharacteristic1
	watchPropertiesChannel chan *dbus.Signal
}

func (d *Device) unwatchProperties() error {
	var err error
	if d.watchPropertiesChannel != nil {
		err = d.client.Unregister(d.watchPropertiesChannel)
		close(d.watchPropertiesChannel)
		d.watchPropertiesChannel = nil
	}

	return err
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

	props, err := c.GetProperties()

	if err != nil {
		return nil, err
	}

	d.Properties = props
	return d.Properties, err
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
func (d *Device) On(name string, fn *emitter.Callback) error {
	switch name {
	case "changed":
		err := d.watchProperties()
		if err != nil {
			return err
		}
		break
	}
	return emitter.On(d.Path+"."+name, fn)
}

//Off unregister callback for event
func (d *Device) Off(name string, cb *emitter.Callback) error {

	var err error

	switch name {
	case "changed":
		err = d.unwatchProperties()
		if err != nil {
			return err
		}
		break
	}

	pattern := d.Path + "." + name
	if name != "*" {
		err = emitter.Off(pattern, cb)
	} else {
		err = emitter.RemoveListeners(pattern, nil)
	}

	return err
}

//Emit an event
func (d *Device) Emit(name string, data interface{}) error {
	return emitter.Emit(d.Path+"."+name, data)
}

//GetService return a GattService
func (d *Device) GetService(path string) *profile.GattService1 {
	return profile.NewGattService1(path, "org.bluez")
}

//GetChar return a GattService
func (d *Device) GetChar(path string) (*profile.GattCharacteristic1, error) {
	return profile.NewGattCharacteristic1(path)
}

//GetAllServicesAndUUID return a list of uuid's with their corresponding services
func (d *Device) GetAllServicesAndUUID() ([]string, error) {

	list, err := d.GetCharsList()
	if err != nil {
		return nil, err
	}

	var deviceFound []string
	var uuidAndService string
	for _, path := range list {

		_, ok := d.chars[path]
		if !ok {
			char, err := profile.NewGattCharacteristic1(string(path))
			if err != nil {
				return nil, err
			}
			d.chars[path] = char
		}

		props := d.chars[path].Properties
		cuuid := strings.ToUpper(props.UUID)
		service := string(props.Service)

		uuidAndService = fmt.Sprint(cuuid, ":", service)
		deviceFound = append(deviceFound, uuidAndService)
	}

	return deviceFound, nil
}

//GetCharByUUID return a GattService by its uuid, return nil if not found
func (d *Device) GetCharByUUID(uuid string) (*profile.GattCharacteristic1, error) {

	uuid = strings.ToUpper(uuid)

	list, err := d.GetCharsList()
	if err != nil {
		return nil, err
	}

	var deviceFound *profile.GattCharacteristic1

	for _, path := range list {

		// use cache
		_, ok := d.chars[path]
		if !ok {
			char, err := profile.NewGattCharacteristic1(string(path))
			if err != nil {
				return nil, err
			}
			d.chars[path] = char
		}

		props := d.chars[path].Properties
		cuuid := strings.ToUpper(props.UUID)

		if cuuid == uuid {
			deviceFound = d.chars[path]
		}
	}

	if deviceFound == nil {
		return nil, errors.New("characteristic not found")
	}

	return deviceFound, nil
}

//GetCharsList return a device characteristics
func (d *Device) GetCharsList() ([]dbus.ObjectPath, error) {

	var chars []dbus.ObjectPath

	manager, err := GetManager()
	if err != nil {
		return nil, err
	}

	list := manager.GetObjects()
	for objpath := range *list {
		path := string(objpath)
		if !strings.HasPrefix(path, d.Path) {
			continue
		}
		charPos := strings.Index(path, "char")
		if charPos == -1 {
			continue
		}
		if strings.Index(path[charPos:], "desc") != -1 {
			continue
		}

		chars = append(chars, objpath)
	}

	return chars, nil
}

//IsConnected check if connected to the device
func (d *Device) IsConnected() bool {

	props, _ := d.GetProperties()

	if props == nil {
		return false
	}

	return props.Connected
}

//Connect to device
func (d *Device) Connect() error {

	c, err := d.GetClient()
	if err != nil {
		return err
	}

	err = c.Connect()
	if err != nil {
		return err
	}
	return nil
}

//Disconnect from a device
func (d *Device) Disconnect() error {
	c, err := d.GetClient()
	if err != nil {
		return err
	}
	c.Disconnect()
	if d.watchPropertiesChannel != nil {
		d.unwatchProperties()
	}
	return nil
}

//Pair a device
func (d *Device) Pair() error {
	c, err := d.GetClient()
	if err != nil {
		return err
	}
	c.Pair()
	return nil
}
