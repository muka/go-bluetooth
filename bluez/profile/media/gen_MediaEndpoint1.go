package media



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var MediaEndpoint1Interface = "org.bluez.MediaEndpoint1"


// NewMediaEndpoint1 create a new instance of MediaEndpoint1
//
// Args:
// - servicePath: unique name
// - objectPath: freely definable
func NewMediaEndpoint1(servicePath string, objectPath dbus.ObjectPath) (*MediaEndpoint1, error) {
	a := new(MediaEndpoint1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: MediaEndpoint1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaEndpoint1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MediaEndpoint1 MediaEndpoint1 hierarchy

type MediaEndpoint1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*MediaEndpoint1Properties
}

// MediaEndpoint1Properties contains the exposed properties of an interface
type MediaEndpoint1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *MediaEndpoint1Properties) Lock() {
	p.lock.Lock()
}

func (p *MediaEndpoint1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *MediaEndpoint1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return MediaEndpoint1 object path
func (a *MediaEndpoint1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return MediaEndpoint1 interface
func (a *MediaEndpoint1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *MediaEndpoint1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a MediaEndpoint1Properties to map
func (a *MediaEndpoint1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaEndpoint1Properties
func (a *MediaEndpoint1Properties) FromMap(props map[string]interface{}) (*MediaEndpoint1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaEndpoint1Properties
func (a *MediaEndpoint1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaEndpoint1Properties, error) {
	s := new(MediaEndpoint1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaEndpoint1) GetProperties() (*MediaEndpoint1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaEndpoint1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaEndpoint1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *MediaEndpoint1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *MediaEndpoint1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *MediaEndpoint1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *MediaEndpoint1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




//SetConfiguration Set configuration for the transport.
func (a *MediaEndpoint1) SetConfiguration(transport dbus.ObjectPath, properties map[string]interface{}) error {
	
	return a.client.Call("SetConfiguration", 0, transport, properties).Store()
	
}

//SelectConfiguration Select preferable configuration from the supported
// capabilities.
// Returns a configuration which can be used to setup
// a transport.
// Note: There is no need to cache the selected
// configuration since on success the configuration is
// send back as parameter of SetConfiguration.
func (a *MediaEndpoint1) SelectConfiguration(capabilities []byte) ([]byte, error) {
	
	var val0 []byte
	err := a.client.Call("SelectConfiguration", 0, capabilities).Store(&val0)
	return val0, err	
}

//ClearConfiguration Clear transport configuration.
func (a *MediaEndpoint1) ClearConfiguration(transport dbus.ObjectPath) error {
	
	return a.client.Call("ClearConfiguration", 0, transport).Store()
	
}

//Release This method gets called when the service daemon
// unregisters the endpoint. An endpoint can use it to do
// cleanup tasks. There is no need to unregister the
// endpoint, because when this method gets called it has
// already been unregistered.
func (a *MediaEndpoint1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

