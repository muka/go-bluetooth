package obex

import (
	log "github.com/sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
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
	client     *bluez.Client
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
func (a *ObexClient1) CreateSession(destination string, options map[string]interface{}) string {
	log.Debugf("CreateSession to %s", destination)
	var sessionPath string
	err := a.client.Call("CreateSession", 0, destination, options).Store(&sessionPath)
	if err != nil {
			panic(err)
	}
	return sessionPath
}

//	Unregister session and abort pending transfers.
//
//	Possible errors:
//		- org.bluez.obex.Error.InvalidArguments
//		- org.bluez.obex.Error.NotAuthorized
//
func (a *ObexClient1) RemoveSession(app dbus.ObjectPath) error {
	return a.client.Call("RemoveSession", 0, app).Store()
}
