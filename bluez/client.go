package bluez

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/util"
)

// NewClient create a new client
func NewClient(config *Config) *Client {
	c := new(Client)
	c.Config = config
	return c
}

// Client implement a DBus client
type Client struct {
	conn       *dbus.Conn
	dbusObject dbus.BusObject
	Config     *Config
}

func (c *Client) isConnected() bool {
	return c.conn != nil
}

//Disconnect from DBus
func (c *Client) Disconnect() {

	// do not disconnect SystemBus
	// as it is a singleton from dbus package
	if c.Config.Bus == SystemBus {
		return
	}

	if c.isConnected() {
		c.conn.Close()
		c.conn = nil
		c.dbusObject = nil
	}
}

// Connect connects to DBus
func (c *Client) Connect() error {
	dbusConn, err := GetConnection(c.Config.Bus)
	if err != nil {
		return err
	}
	c.conn = dbusConn
	c.dbusObject = c.conn.Object(c.Config.Name, dbus.ObjectPath(c.Config.Path))
	return nil
}

// Call a DBus method
func (c *Client) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return &dbus.Call{
				Err: err,
			}
		}
	}

	methodPath := fmt.Sprint(c.Config.Iface, ".", method)
	return c.dbusObject.Call(methodPath, flags, args...)
}

//GetConnection returns the Dbus connection
func (c *Client) GetConnection() *dbus.Conn {
	return c.conn
}

//GetDbusObject returns the Dbus object
func (c *Client) GetDbusObject() dbus.BusObject {
	return c.dbusObject
}

//GetProperty return a property value
func (c *Client) GetProperty(p string) (dbus.Variant, error) {
	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return dbus.Variant{}, err
		}
	}
	return c.dbusObject.GetProperty(c.Config.Iface + "." + p)
}

//SetProperty set a property value
func (c *Client) SetProperty(p string, v interface{}) error {
	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}
	return c.dbusObject.Call("org.freedesktop.DBus.Properties.Set", 0, c.Config.Iface, p, dbus.MakeVariant(v)).Store()
}

//GetProperties load all the properties for an interface
func (c *Client) GetProperties(props interface{}) error {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}

	result := make(map[string]dbus.Variant)
	err := c.dbusObject.Call("org.freedesktop.DBus.Properties.GetAll", 0, c.Config.Iface).Store(&result)
	if err != nil {
		return fmt.Errorf("Properties.GetAll %s: %s", c.Config.Iface, err)
	}

	err = util.MapToStruct(props, result)
	if err != nil {
		return fmt.Errorf("MapToStruct: %s", err)
	}

	return nil
}

func getMatchString(path dbus.ObjectPath, iface string) string {
	return fmt.Sprintf("type='signal',interface='%s',path='%s'", iface, path)
}

//Register for signals
func (c *Client) Register(path dbus.ObjectPath, iface string) (chan *dbus.Signal, error) {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return nil, err
		}
	}

	matchstr := getMatchString(path, iface)
	c.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchstr)

	channel := make(chan *dbus.Signal, 1)
	c.conn.Signal(channel)

	return channel, nil
}

//Unregister for signals
func (c *Client) Unregister(path dbus.ObjectPath, iface string, signal chan *dbus.Signal) error {
	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}
	matchstr := getMatchString(path, iface)
	c.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, matchstr)
	if signal != nil {
		c.conn.RemoveSignal(signal)
	}

	return nil
}

// Emit
func (c *Client) Emit(path dbus.ObjectPath, name string, values ...interface{}) error {
	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}
	return c.conn.Emit(path, name, values...)
}
