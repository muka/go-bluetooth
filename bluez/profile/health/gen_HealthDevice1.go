// WARNING: generated code, do not edit!
// Copyright © 2019 luca capra
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package health



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var HealthDevice1Interface = "org.bluez.HealthDevice1"


// NewHealthDevice1 create a new instance of HealthDevice1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewHealthDevice1(objectPath dbus.ObjectPath) (*HealthDevice1, error) {
	a := new(HealthDevice1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: HealthDevice1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(HealthDevice1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// HealthDevice1 HealthDevice hierarchy

type HealthDevice1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*HealthDevice1Properties
}

// HealthDevice1Properties contains the exposed properties of an interface
type HealthDevice1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// MainChannel The first reliable channel opened. It is needed by
  // upper applications in order to send specific protocol
  // data units. The first reliable can change after a
  // reconnection.
	MainChannel dbus.ObjectPath

}

func (p *HealthDevice1Properties) Lock() {
	p.lock.Lock()
}

func (p *HealthDevice1Properties) Unlock() {
	p.lock.Unlock()
}


// SetMainChannel set MainChannel value
func (a *HealthDevice1) SetMainChannel(v dbus.ObjectPath) error {
	return a.SetProperty("MainChannel", v)
}

// GetMainChannel get MainChannel value
func (a *HealthDevice1) GetMainChannel() (dbus.ObjectPath, error) {
	v, err := a.GetProperty("MainChannel")
	if err != nil {
		return dbus.ObjectPath(""), err
	}
	return v.Value().(dbus.ObjectPath), nil
}


// Close the connection
func (a *HealthDevice1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return HealthDevice1 object path
func (a *HealthDevice1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return HealthDevice1 interface
func (a *HealthDevice1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *HealthDevice1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}


// ToMap convert a HealthDevice1Properties to map
func (a *HealthDevice1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an HealthDevice1Properties
func (a *HealthDevice1Properties) FromMap(props map[string]interface{}) (*HealthDevice1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an HealthDevice1Properties
func (a *HealthDevice1Properties) FromDBusMap(props map[string]dbus.Variant) (*HealthDevice1Properties, error) {
	s := new(HealthDevice1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *HealthDevice1) GetProperties() (*HealthDevice1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *HealthDevice1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *HealthDevice1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *HealthDevice1) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *HealthDevice1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *HealthDevice1) WatchProperties() (chan *bluez.PropertyChanged, error) {

	channel, err := a.client.Register(a.Path(), a.Interface())
	if err != nil {
		return nil, err
	}

	ch := make(chan *bluez.PropertyChanged)

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
			if sig.Path != a.Path() {
				continue
			}

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)

			for field, val := range changes {

				// updates [*]Properties struct when a property change
				s := reflect.ValueOf(a.Properties).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						a.Properties.Lock()
						f.Set(x)
						a.Properties.Unlock()
					}
				}

				propChanged := &bluez.PropertyChanged{
					Interface: iface,
					Name:      field,
					Value:     val.Value(),
				}
				ch <- propChanged
			}

		}
	})()

	return ch, nil
}

func (a *HealthDevice1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




//Echo Sends an echo petition to the remote service. Returns
// True if response matches with the buffer sent. If some
// error is detected False value is returned.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.OutOfRange
func (a *HealthDevice1) Echo() (bool, error) {
	
	var val0 bool
	err := a.client.Call("Echo", 0, ).Store(&val0)
	return val0, err	
}

//CreateChannel Creates a new data channel.  The configuration should
// indicate the channel quality of service using one of
// this values "reliable", "streaming", "any".
// Returns the object path that identifies the data
// channel that is already connected.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.HealthError
func (a *HealthDevice1) CreateChannel(application dbus.ObjectPath, configuration string) (dbus.ObjectPath, error) {
	
	var val0 dbus.ObjectPath
	err := a.client.Call("CreateChannel", 0, application, configuration).Store(&val0)
	return val0, err	
}

//DestroyChannel Destroys the data channel object. Only the creator of
// the channel or the creator of the HealthApplication
// that received the data channel will be able to destroy
// it.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.NotFound
// org.bluez.Error.NotAllowed
func (a *HealthDevice1) DestroyChannel(channel dbus.ObjectPath) error {
	
	return a.client.Call("DestroyChannel", 0, channel).Store()
	
}

