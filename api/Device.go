package api

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/util"
)

var deviceRegistry = new(sync.Map)

// NewDevice creates a new Device
func NewDevice(path string) *Device {

	if device, ok := deviceRegistry.Load(path); ok {
		return device.(*Device)
	}

	d := new(Device)
	d.Path = path
	d.client = profile.NewDevice1(path)

	d.client.GetProperties()
	d.Properties = d.client.Properties
	d.chars = make(map[dbus.ObjectPath]*profile.GattCharacteristic1, 0)
	d.descr = make(map[dbus.ObjectPath]*profile.GattDescriptor1, 0)

	deviceRegistry.Store(path, d)

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

	if _, ok := deviceRegistry.Load(d.Path); ok {
		deviceRegistry.Delete(d.Path)
	}

	return nil
}

// FlushDevices clears removes device and clears it from cached devices
func FlushDevices(adapterID string) error {
	adapter, err := GetAdapter(adapterID)
	if err != nil {
		return err
	}

	devices, err := GetDevices()
	if err != nil {
		return err
	}
	for _, dev := range devices {
		err := ClearDevice(&dev)
		if err != nil {
			return fmt.Errorf("FlushDevices.ClearDevice %s", err)
		}
		err = adapter.RemoveDevice(dev.Path)
		if err != nil {
			return fmt.Errorf("FlushDevices.RemoveDevice %s", err)
		}
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
	d.descr = make(map[dbus.ObjectPath]*profile.GattDescriptor1, 0)

	return d, nil
}

func (d *Device) watchProperties() error {
	d.lock.RLock()
	if d.watchPropertiesChannel != nil {
		d.lock.RUnlock()
		d.unwatchProperties()
	} else {
		d.lock.RUnlock()
	}

	channel, err := d.client.Register()
	if err != nil {
		return err
	}
	d.lock.Lock()
	d.watchPropertiesChannel = channel
	d.lock.Unlock()

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
			if fmt.Sprint(sig.Path) != d.Path {
				continue
			}

			// for i := 0; i < len(sig.Body); i++ {
			// log.Printf("%s -> %s\n", reflect.TypeOf(sig.Body[i]), sig.Body[i])
			// }

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)
			propertyChangedEvents := make([]PropertyChangedEvent, 0)
			for field, val := range changes {

				// updates [*]Properties struct
				d.lock.RLock()
				props := d.Properties
				d.lock.RUnlock()

				s := reflect.ValueOf(props).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						props.Lock.Lock()
						f.Set(x)
						props.Lock.Unlock()
					}
				}

				propChanged := PropertyChangedEvent{string(iface), field, val.Value(), props, d}
				propertyChangedEvents = append(propertyChangedEvents, propChanged)
			}

			for _, propChanged := range propertyChangedEvents {
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
	descr                  map[dbus.ObjectPath]*profile.GattDescriptor1
	watchPropertiesChannel chan *dbus.Signal
	lock                   sync.RWMutex
}

func (d *Device) unwatchProperties() error {
	var err error
	d.lock.Lock()
	defer d.lock.Unlock()
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

	d.lock.Lock()
	d.Properties = props
	defer d.lock.Unlock()

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
	devices, err := d.GetCharsByUUID(uuid)
	if len(devices) > 0 {
		return devices[0], err
	}
	return nil, err
}

//GetCharsByUUID returns all characteristics that match the given UUID.
func (d *Device) GetCharsByUUID(uuid string) ([]*profile.GattCharacteristic1, error) {
	uuid = strings.ToUpper(uuid)

	list, err := d.GetCharsList()
	if err != nil {
		return nil, err
	}

	var charsFound []*profile.GattCharacteristic1

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
			charsFound = append(charsFound, d.chars[path])
		}
	}

	if len(charsFound) == 0 {
		return nil, errors.New("characteristic not found")
	}

	return charsFound, nil
}

//GetCharsList return a device characteristics
func (d *Device) GetCharsList() ([]dbus.ObjectPath, error) {

	var chars []dbus.ObjectPath

	manager, err := GetManager()
	if err != nil {
		return nil, err
	}

	list := manager.GetObjects()
	list.Range(func(objpath, value interface{}) bool {
		path := string(objpath.(dbus.ObjectPath))
		if !strings.HasPrefix(path, d.Path) {
			return true
		}
		charPos := strings.Index(path, "char")
		if charPos == -1 {
			return true
		}
		if strings.Index(path[charPos:], "desc") != -1 {
			return true
		}

		chars = append(chars, objpath.(dbus.ObjectPath))
		return true
	})

	return chars, nil
}

//GetDescriptorList returns all descriptors
func (d *Device) GetDescriptorList() ([]dbus.ObjectPath, error) {
	var descr []dbus.ObjectPath

	manager, err := GetManager()
	if err != nil {
		return nil, err
	}

	list := manager.GetObjects()
	list.Range(func(objpath, value interface{}) bool {
		path := string(objpath.(dbus.ObjectPath))
		if !strings.HasPrefix(path, d.Path) {
			return true
		}
		charPos := strings.Index(path, "char")
		if charPos == -1 {
			return true
		}
		if strings.Index(path[charPos:], "desc") == -1 {
			return true
		}

		descr = append(descr, objpath.(dbus.ObjectPath))
		return true
	})

	return descr, nil
}

//GetDescriptors returns all descriptors for a given characteristic
func (d *Device) GetDescriptors(char *profile.GattCharacteristic1) ([]*profile.GattDescriptor1, error) {
	descrPaths, err := d.GetDescriptorList()
	if err != nil {
		return nil, err
	}

	var descrFound []*profile.GattDescriptor1

	for _, path := range descrPaths {
		_, ok := d.descr[path]
		if !ok {
			descr, err := profile.NewGattDescriptor1(string(path))
			if err != nil {
				return nil, err
			}
			d.descr[path] = descr
		}

		if dbus.ObjectPath(char.Path) == d.descr[path].Properties.Characteristic {
			descrFound = append(descrFound, d.descr[path])
		}
	}

	if len(descrFound) == 0 {
		return nil, errors.New("descriptors not found")
	}

	return descrFound, nil
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

	d.lock.RLock()
	if d.watchPropertiesChannel != nil {
		d.lock.RUnlock()
		d.unwatchProperties()
	} else {
		d.lock.RUnlock()
	}
	return nil
}

// Pair a device
func (d *Device) Pair() error {
	c, err := d.GetClient()
	if err != nil {
		return err
	}
	return c.Pair()
}
