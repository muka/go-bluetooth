package bluez

import (
	"log"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/util"
	"github.com/muka/device-manager/client"
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

// Close the connection
func (d *Device1) Close() {
	d.client.Disconnect()
}

//GetProperties load all available properties
func (d *Device1) GetProperties() (*Device1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//CancelParing stop the pairing process
func (d *Device1) CancelParing() error {
	return d.client.Call("CancelParing", 0).Store()
}

//Connect to the device
func (d *Device1) Connect() error {
	return d.client.Call("Connect", 0).Store()
}

//ConnectProfile connect to the specific profile
func (d *Device1) ConnectProfile(uuid string) error {
	return d.client.Call("ConnectProfile", 0, uuid).Store()
}

//Disconnect from the device
func (d *Device1) Disconnect() error {
	return d.client.Call("Disconnect", 0).Store()
}

//DisconnectProfile from the device
func (d *Device1) DisconnectProfile(uuid string) error {
	return d.client.Call("DisconnectProfile", 0, uuid).Store()
}

//Pair with the device
func (d *Device1) Pair() error {
	return d.client.Call("Pair", 0).Store()
}
