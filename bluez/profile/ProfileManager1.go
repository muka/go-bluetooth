package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

// NewProfileManager1 create a new ProfileManager1 client
func NewProfileManager1(hostID string) *ProfileManager1 {
	a := new(ProfileManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.ProfileManager1",
			Path:  "/org/bluez",
			Bus:   bluez.SystemBus,
		},
	)
	return a
}

// ProfileManager1 client
type ProfileManager1 struct {
	client *bluez.Client
}

// Close the connection
func (a *ProfileManager1) Close() {
	a.client.Disconnect()
}

//RegisterProfile add a new Profile for an UUID
func (a *ProfileManager1) RegisterProfile(profile string, UUID string, options map[string]interface{}) error {
	return a.client.Call("RegisterProfile", 0, dbus.ObjectPath(profile), UUID, options).Store()
}

//UnregisterProfile add a new Profile for an UUID
func (a *ProfileManager1) UnregisterProfile(profile string) error {
	return a.client.Call("UnregisterProfile", 0, dbus.ObjectPath(profile)).Store()
}
