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

package network



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Network1Interface = "org.bluez.Network1"


// NewNetwork1 create a new instance of Network1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewNetwork1(objectPath dbus.ObjectPath) (*Network1, error) {
	a := new(Network1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Network1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Network1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Network1 Network hierarchy

type Network1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Network1Properties
}

// Network1Properties contains the exposed properties of an interface
type Network1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Connected Indicates if the device is connected.
	Connected bool

	// Interface Indicates the network interface name when available.
	Interface string

	// UUID Indicates the connection role when available.
	UUID string

}

func (p *Network1Properties) Lock() {
	p.lock.Lock()
}

func (p *Network1Properties) Unlock() {
	p.lock.Unlock()
}


// SetConnected set Connected value
func (a *Network1) SetConnected(v bool) error {
	return a.SetProperty("Connected", v)
}

// GetConnected get Connected value
func (a *Network1) GetConnected() (bool, error) {
	v, err := a.GetProperty("Connected")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetInterface set Interface value
func (a *Network1) SetInterface(v string) error {
	return a.SetProperty("Interface", v)
}

// GetInterface get Interface value
func (a *Network1) GetInterface() (string, error) {
	v, err := a.GetProperty("Interface")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetUUID set UUID value
func (a *Network1) SetUUID(v string) error {
	return a.SetProperty("UUID", v)
}

// GetUUID get UUID value
func (a *Network1) GetUUID() (string, error) {
	v, err := a.GetProperty("UUID")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}


// Close the connection
func (a *Network1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Network1 object path
func (a *Network1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return Network1 interface
func (a *Network1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Network1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a Network1Properties to map
func (a *Network1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Network1Properties
func (a *Network1Properties) FromMap(props map[string]interface{}) (*Network1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Network1Properties
func (a *Network1Properties) FromDBusMap(props map[string]dbus.Variant) (*Network1Properties, error) {
	s := new(Network1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Network1) GetProperties() (*Network1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Network1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Network1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Network1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Network1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Network1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *Network1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




//Connect Connect to the network device and return the network
// interface name. Examples of the interface name are
// bnep0, bnep1 etc.
// uuid can be either one of "gn", "panu" or "nap" (case
// insensitive) or a traditional string representation of
// UUID or a hexadecimal number.
// The connection will be closed and network device
func (a *Network1) Connect(uuid string) (string, error) {
	
	var val0 string
	err := a.client.Call("Connect", 0, uuid).Store(&val0)
	return val0, err	
}

//Disconnect Disconnect from the network device.
// To abort a connection attempt in case of errors or
// timeouts in the client it is fine to call this method.
// Possible errors: org.bluez.Error.Failed
func (a *Network1) Disconnect() error {
	
	return a.client.Call("Disconnect", 0, ).Store()
	
}

