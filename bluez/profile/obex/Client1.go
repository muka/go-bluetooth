package obex

import (
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

// TODO: https://github.com/blueman-project/blueman/issues/218#issuecomment-89315974
// NewObexClient1 create a new ObexClient1 client
func NewObexClient1() *ObexClient1 {
	a := new(ObexClient1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: "org.bluez.obex.Client1",
			Path:  "/org/bluez/obex",
			Bus:   bluez.SessionBus,
		},
	)
	return a
}

// ObexClient1 client
type ObexClient1 struct {
	client *bluez.Client
}

// Close the connection
func (a *ObexClient1) Close() {
	a.client.Disconnect()
}

//	Create a new OBEX session for the given remote address.
//
//	The last parameter is a dictionary to hold optional or
//	type-specific parameters. Typical parameters that can
//	be set in this dictionary include the following:
//
//		string "Target" : type of session to be created
//		string "Source" : local address to be used
//		byte "Channel"
//
//	The currently supported targets are the following:
//
//	- "ftp"
//	- "map"
//	- "opp"
//	- "pbap"
//	- "sync"
//
//	Possible errors:
// 		- org.bluez.obex.Error.InvalidArguments
//		- org.bluez.obex.Error.Failed
//
// TODO: Use ObexSession1 struct instead of generic map for options
func (a *ObexClient1) CreateSession(destination string, options map[string]interface{}) (string, error) {
	log.Debugf("CreateSession to %s", destination)
	var sessionPath string
	err := a.client.Call("CreateSession", 0, destination, options).Store(&sessionPath)
	return sessionPath, err
}

//	Unregister session and abort pending transfers.
//
//	Possible errors:
//		- org.bluez.obex.Error.InvalidArguments
//		- org.bluez.obex.Error.NotAuthorized
//
func (a *ObexClient1) RemoveSession(session string) error {
	return a.client.Call("RemoveSession", 0, dbus.ObjectPath(session)).Store()
}
