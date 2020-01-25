package profile



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/muka/go-bluetooth/util"
   "github.com/muka/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var ProfileManager1Interface = "org.bluez.ProfileManager1"


// NewProfileManager1 create a new instance of ProfileManager1
//
// Args:

func NewProfileManager1() (*ProfileManager1, error) {
	a := new(ProfileManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: ProfileManager1Interface,
			Path:  dbus.ObjectPath("/org/bluez"),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(ProfileManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
ProfileManager1 Profile Manager hierarchy

*/
type ProfileManager1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*ProfileManager1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// ProfileManager1Properties contains the exposed properties of an interface
type ProfileManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *ProfileManager1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *ProfileManager1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *ProfileManager1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return ProfileManager1 object path
func (a *ProfileManager1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return ProfileManager1 dbus client
func (a *ProfileManager1) Client() *bluez.Client {
	return a.client
}

// Interface return ProfileManager1 interface
func (a *ProfileManager1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *ProfileManager1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a ProfileManager1Properties to map
func (a *ProfileManager1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an ProfileManager1Properties
func (a *ProfileManager1Properties) FromMap(props map[string]interface{}) (*ProfileManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an ProfileManager1Properties
func (a *ProfileManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*ProfileManager1Properties, error) {
	s := new(ProfileManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *ProfileManager1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *ProfileManager1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *ProfileManager1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *ProfileManager1) GetProperties() (*ProfileManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *ProfileManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *ProfileManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *ProfileManager1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *ProfileManager1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *ProfileManager1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *ProfileManager1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
RegisterProfile 
			This registers a profile implementation.

			If an application disconnects from the bus all
			its registered profiles will be removed.

			HFP HS UUID: 0000111e-0000-1000-8000-00805f9b34fb

				Default RFCOMM channel is 6. And this requires
				authentication.

			Available options:

				string Name

					Human readable name for the profile

				string Service


*/
func (a *ProfileManager1) RegisterProfile(profile dbus.ObjectPath, uuid string, options map[string]interface{}) error {
	
	return a.client.Call("RegisterProfile", 0, profile, uuid, options).Store()
	
}

/*
UnregisterProfile 
			This unregisters the profile that has been previously
			registered. The object path parameter must match the
			same value that has been used on registration.

			Possible errors: org.bluez.Error.DoesNotExist



*/
func (a *ProfileManager1) UnregisterProfile(profile dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterProfile", 0, profile).Store()
	
}

