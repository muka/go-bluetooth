package gatt



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var GattDescriptor1Interface = "org.bluez.GattDescriptor1"


// NewGattDescriptor1 create a new instance of GattDescriptor1
//
// Args:
// - objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/serviceXX/charYYYY/descriptorZZZ
func NewGattDescriptor1(objectPath dbus.ObjectPath) (*GattDescriptor1, error) {
	a := new(GattDescriptor1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattDescriptor1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattDescriptor1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
GattDescriptor1 Characteristic Descriptors hierarchy
Local or remote GATT characteristic descriptors hierarchy.
*/
type GattDescriptor1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*GattDescriptor1Properties
}

// GattDescriptor1Properties contains the exposed properties of an interface
type GattDescriptor1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	/*
	Value The cached value of the descriptor. This property
  gets updated only after a successful read request, upon
  which a PropertiesChanged signal will be emitted.
	*/
	Value []byte `dbus:"emit"`

	/*
	Flags Defines how the descriptor value can be used.
  Possible values:
  "read"
  "write"
  "encrypt-read"
  "encrypt-write"
  "encrypt-authenticated-read"
  "encrypt-authenticated-write"
  "secure-read" (Server Only)
  "secure-write" (Server Only)
  "authorize"
	*/
	Flags []string

	/*
	UUID 128-bit descriptor UUID.
	*/
	UUID string

	/*
	Characteristic Object path of the GATT characteristic the descriptor
  belongs to.
	*/
	Characteristic dbus.ObjectPath

}

//Lock access to properties
func (p *GattDescriptor1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *GattDescriptor1Properties) Unlock() {
	p.lock.Unlock()
}


// SetValue set Value value
func (a *GattDescriptor1) SetValue(v []byte) error {
	return a.SetProperty("Value", v)
}

// GetValue get Value value
func (a *GattDescriptor1) GetValue() ([]byte, error) {
	v, err := a.GetProperty("Value")
	if err != nil {
		return []byte{}, err
	}
	return v.Value().([]byte), nil
}

// SetFlags set Flags value
func (a *GattDescriptor1) SetFlags(v []string) error {
	return a.SetProperty("Flags", v)
}

// GetFlags get Flags value
func (a *GattDescriptor1) GetFlags() ([]string, error) {
	v, err := a.GetProperty("Flags")
	if err != nil {
		return []string{}, err
	}
	return v.Value().([]string), nil
}

// SetUUID set UUID value
func (a *GattDescriptor1) SetUUID(v string) error {
	return a.SetProperty("UUID", v)
}

// GetUUID get UUID value
func (a *GattDescriptor1) GetUUID() (string, error) {
	v, err := a.GetProperty("UUID")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetCharacteristic set Characteristic value
func (a *GattDescriptor1) SetCharacteristic(v dbus.ObjectPath) error {
	return a.SetProperty("Characteristic", v)
}

// GetCharacteristic get Characteristic value
func (a *GattDescriptor1) GetCharacteristic() (dbus.ObjectPath, error) {
	v, err := a.GetProperty("Characteristic")
	if err != nil {
		return dbus.ObjectPath(""), err
	}
	return v.Value().(dbus.ObjectPath), nil
}


// Close the connection
func (a *GattDescriptor1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return GattDescriptor1 object path
func (a *GattDescriptor1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return GattDescriptor1 interface
func (a *GattDescriptor1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *GattDescriptor1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a GattDescriptor1Properties to map
func (a *GattDescriptor1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an GattDescriptor1Properties
func (a *GattDescriptor1Properties) FromMap(props map[string]interface{}) (*GattDescriptor1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattDescriptor1Properties
func (a *GattDescriptor1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattDescriptor1Properties, error) {
	s := new(GattDescriptor1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *GattDescriptor1) GetProperties() (*GattDescriptor1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattDescriptor1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattDescriptor1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *GattDescriptor1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *GattDescriptor1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *GattDescriptor1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *GattDescriptor1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




/*
ReadValue Issues a request to read the value of the
characteristic and returns the value if the
operation was successful.
Possible options: "offset": Start offset
"device": Device path (Server only)
"link": Link type (Server only)
Possible Errors: org.bluez.Error.Failed
org.bluez.Error.InProgress
org.bluez.Error.NotPermitted
org.bluez.Error.NotAuthorized
org.bluez.Error.NotSupported
*/
func (a *GattDescriptor1) ReadValue(flags map[string]interface{}) ([]byte, error) {
	
	var val0 []byte
	err := a.client.Call("ReadValue", 0, flags).Store(&val0)
	return val0, err	
}

/*
WriteValue Issues a request to write the value of the
characteristic.
Possible options: "offset": Start offset
"device": Device path (Server only)
"link": Link type (Server only)
"prepare-authorize": boolean Is prepare
authorization
request
Possible Errors: org.bluez.Error.Failed
org.bluez.Error.InProgress
org.bluez.Error.NotPermitted
org.bluez.Error.InvalidValueLength
org.bluez.Error.NotAuthorized
org.bluez.Error.NotSupported
*/
func (a *GattDescriptor1) WriteValue(value []byte, flags map[string]interface{}) error {
	
	return a.client.Call("WriteValue", 0, value, flags).Store()
	
}

