package thermometer



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var ThermometerWatcher1Interface = "org.bluez.ThermometerWatcher1"


// NewThermometerWatcher1 create a new instance of ThermometerWatcher1
//
// Args:
// - servicePath: unique name
// - objectPath: freely definable
func NewThermometerWatcher1(servicePath string, objectPath dbus.ObjectPath) (*ThermometerWatcher1, error) {
	a := new(ThermometerWatcher1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: ThermometerWatcher1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(ThermometerWatcher1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// ThermometerWatcher1 Health Thermometer Watcher hierarchy

type ThermometerWatcher1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*ThermometerWatcher1Properties
}

// ThermometerWatcher1Properties contains the exposed properties of an interface
type ThermometerWatcher1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *ThermometerWatcher1Properties) Lock() {
	p.lock.Lock()
}

func (p *ThermometerWatcher1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *ThermometerWatcher1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return ThermometerWatcher1 object path
func (a *ThermometerWatcher1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return ThermometerWatcher1 interface
func (a *ThermometerWatcher1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *ThermometerWatcher1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a ThermometerWatcher1Properties to map
func (a *ThermometerWatcher1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an ThermometerWatcher1Properties
func (a *ThermometerWatcher1Properties) FromMap(props map[string]interface{}) (*ThermometerWatcher1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an ThermometerWatcher1Properties
func (a *ThermometerWatcher1Properties) FromDBusMap(props map[string]dbus.Variant) (*ThermometerWatcher1Properties, error) {
	s := new(ThermometerWatcher1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *ThermometerWatcher1) GetProperties() (*ThermometerWatcher1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *ThermometerWatcher1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *ThermometerWatcher1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *ThermometerWatcher1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *ThermometerWatcher1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *ThermometerWatcher1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *ThermometerWatcher1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




//MeasurementReceived This callback gets called when a measurement has been
// scanned in the thermometer.
// Measurement:
// int16 Exponent:
// int32 Mantissa:
// Exponent and Mantissa values as
// extracted from float value defined by
// IEEE-11073-20601.
func (a *ThermometerWatcher1) MeasurementReceived(measurement map[string]interface{}) error {
	
	return a.client.Call("MeasurementReceived", 0, measurement).Store()
	
}

