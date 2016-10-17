package bluez

import (
	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/util"
	"github.com/muka/device-manager/client"
	"log"
)

// NewDevice1 create a new Device1 client
func NewDevice1(path string) *Device1 {
	a := new(Device1)
	a.client = client.NewClient(
		&client.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.Device1",
			Path:  path,
			Bus:   client.SystemBus,
		},
	)
	a.logger = util.NewLogger(path)
	a.Properties = new(Device1Properties)
	return a
}

// Device1 client
type Device1 struct {
	client     *client.Client
	logger     *log.Logger
	Properties *Device1Properties
}

// Device1Properties exposed properties for Device1
type Device1Properties struct {
	UUIDs         []string
	Blocked       bool
	Connected     bool
	LegacyPairing bool
	Paired        bool
	Trusted       bool
	RSSI          int16
	Adapter       dbus.ObjectPath
	Address       string
	Alias         string
	Icon          string
	Modalias      string
	Name          string
	Appearance    uint16
	Class         uint32
}

//GetProperties load all available properties
func (d *Device1) GetProperties() (*Device1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//CancelParing stop the pairing process
func (d *Device1) CancelParing() error {
	return d.client.Call("CancelParing", dbus.FlagNoReplyExpected).Store()
}

//Connect to the device
func (d *Device1) Connect() error {
	return d.client.Call("Connect", dbus.FlagNoReplyExpected).Store()
}

//ConnectProfile connect to the specific profile
func (d *Device1) ConnectProfile(uuid string) error {
	return d.client.Call("ConnectProfile", dbus.FlagNoReplyExpected, uuid).Store()
}

//Disconnect from the device
func (d *Device1) Disconnect() error {
	return d.client.Call("Disconnect", dbus.FlagNoReplyExpected).Store()
}

//DisconnectProfile from the device
func (d *Device1) DisconnectProfile(uuid string) error {
	return d.client.Call("DisconnectProfile", dbus.FlagNoReplyExpected, uuid).Store()
}

//Pair with the device
func (d *Device1) Pair() error {
	return d.client.Call("Pair", dbus.FlagNoReplyExpected).Store()
}
