// Code generated DO NOT EDIT

package profile



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/godbus/dbus/v5"
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




/*
RegisterProfile 			This registers a profile implementation.
			If an application disconnects from the bus all
			its registered profiles will be removed.
			HFP HS UUID: 0000111e-0000-1000-8000-00805f9b34fb
				Default RFCOMM channel is 6. And this requires
				authentication.
			Available options:
				string Name
					Human readable name for the profile
				string Service
					The primary service class UUID
					(if different from the actual
					 profile UUID)
				string Role
					For asymmetric profiles that do not
					have UUIDs available to uniquely
					identify each side this
					parameter allows specifying the
					precise local role.
					Possible values: "client", "server"
				uint16 Channel
					RFCOMM channel number that is used
					for client and server UUIDs.
					If applicable it will be used in the
					SDP record as well.
				uint16 PSM
					PSM number that is used for client
					and server UUIDs.
					If applicable it will be used in the
					SDP record as well.
				boolean RequireAuthentication
					Pairing is required before connections
					will be established. No devices will
					be connected if not paired.
				boolean RequireAuthorization
					Request authorization before any
					connection will be established.
				boolean AutoConnect
					In case of a client UUID this will
					force connection of the RFCOMM or
					L2CAP channels when a remote device
					is connected.
				string ServiceRecord
					Provide a manual SDP record.
				uint16 Version
					Profile version (for SDP record)
				uint16 Features
					Profile features (for SDP record)
			Possible errors: org.bluez.Error.InvalidArguments
			                 org.bluez.Error.AlreadyExists

*/
func (a *ProfileManager1) RegisterProfile(profile dbus.ObjectPath, uuid string, options map[string]interface{}) error {
	
	return a.client.Call("RegisterProfile", 0, profile, uuid, options).Store()
	
}

/*
UnregisterProfile 			This unregisters the profile that has been previously
			registered. The object path parameter must match the
			same value that has been used on registration.
			Possible errors: org.bluez.Error.DoesNotExist

*/
func (a *ProfileManager1) UnregisterProfile(profile dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterProfile", 0, profile).Store()
	
}

