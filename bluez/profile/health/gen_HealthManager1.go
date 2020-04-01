// Code generated DO NOT EDIT

package health



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/muka/go-bluetooth/util"
   "github.com/muka/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var HealthManager1Interface = "org.bluez.HealthManager1"


// NewHealthManager1 create a new instance of HealthManager1
//
// Args:

func NewHealthManager1() (*HealthManager1, error) {
	a := new(HealthManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: HealthManager1Interface,
			Path:  dbus.ObjectPath("/org/bluez/"),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(HealthManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
HealthManager1 HealthManager hierarchy

*/
type HealthManager1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*HealthManager1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// HealthManager1Properties contains the exposed properties of an interface
type HealthManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *HealthManager1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *HealthManager1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *HealthManager1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return HealthManager1 object path
func (a *HealthManager1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return HealthManager1 dbus client
func (a *HealthManager1) Client() *bluez.Client {
	return a.client
}

// Interface return HealthManager1 interface
func (a *HealthManager1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *HealthManager1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a HealthManager1Properties to map
func (a *HealthManager1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an HealthManager1Properties
func (a *HealthManager1Properties) FromMap(props map[string]interface{}) (*HealthManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an HealthManager1Properties
func (a *HealthManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*HealthManager1Properties, error) {
	s := new(HealthManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *HealthManager1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *HealthManager1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *HealthManager1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *HealthManager1) GetProperties() (*HealthManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *HealthManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *HealthManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *HealthManager1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *HealthManager1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *HealthManager1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *HealthManager1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
CreateApplication 
			Returns the path of the new registered application.
			Application will be closed by the call or implicitly
			when the programs leaves the bus.

			config:
				uint16 DataType:

					Mandatory

				string Role:

					Mandatory. Possible values: "source",
									"sink"

				string Description:

					Optional

				ChannelType:

					Optional, just for sources. Possible
					values: "reliable", "streaming"

			Possible Errors: org.bluez.Error.InvalidArguments


*/
func (a *HealthManager1) CreateApplication(config map[string]interface{}) (dbus.ObjectPath, error) {
	
	var val0 dbus.ObjectPath
	err := a.client.Call("CreateApplication", 0, config).Store(&val0)
	return val0, err	
}

/*
DestroyApplication 
			Closes the HDP application identified by the object
			path. Also application will be closed if the process
			that started it leaves the bus. Only the creator of the
			application will be able to destroy it.

			Possible errors: org.bluez.Error.InvalidArguments
					 org.bluez.Error.NotFound
					 org.bluez.Error.NotAllowed



*/
func (a *HealthManager1) DestroyApplication(application dbus.ObjectPath) error {
	
	return a.client.Call("DestroyApplication", 0, application).Store()
	
}

