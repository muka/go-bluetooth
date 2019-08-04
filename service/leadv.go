package service

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
)

// NewLEAdvertisement1 create a new LEAdvertisement1 client
func NewLEAdvertisement1(config *LEAdvertisement1Config, props *advertising.LEAdvertisement1Properties) (*LEAdvertisement1, error) {
	s := &LEAdvertisement1{
		config:     config,
		properties: props,
	}
	return s, nil
}

func NewLEAdvertisement1Config(objectPath dbus.ObjectPath, conn *dbus.Conn) (*LEAdvertisement1Config, error) {
	if conn == nil {
		sysconn, err := dbus.SystemBus()
		if err != nil {
			return nil, err
		}
		conn = sysconn
	}
	return &LEAdvertisement1Config{objectPath, conn}, nil
}

//LEAdvertisement1Config LEAdvertisement1 configuration
type LEAdvertisement1Config struct {
	objectPath dbus.ObjectPath
	conn       *dbus.Conn
}

// LEAdvertisement1 client
type LEAdvertisement1 struct {
	config     *LEAdvertisement1Config
	properties *advertising.LEAdvertisement1Properties
}

//Interface return the dbus interface name
func (s *LEAdvertisement1) Interface() string {
	return advertising.LEAdvertisement1Interface
}

//Path return the object path
func (s *LEAdvertisement1) Path() dbus.ObjectPath {
	return s.config.objectPath
}

//Conn return the object connection
func (s *LEAdvertisement1) Conn() *dbus.Conn {
	return s.config.conn
}

//Properties return the properties of the service
func (s *LEAdvertisement1) Properties() bluez.Properties {
	return s.properties
}

// Release This method gets called when the service daemon
// removes the Advertisement. A client can use it to do
// cleanup tasks. There is no need to call
// UnregisterAdvertisement because when this method gets
// called it has already been unregistered.
func (s *LEAdvertisement1) Release() error {
	return s.Conn().BusObject().Call(s.Interface()+".Release", 0).Store()
}

//Expose the char to dbus
func (s *LEAdvertisement1) Expose() error {
	return ExposeService(s)
}
