package bluez

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/util"
)

// NewClient create a new client
func NewClient(config *Config) *Client {
	dbg("Create new client: %v", config)
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
	if c.isConnected() {
		c.conn.Close()
		c.conn = nil
		c.dbusObject = nil
		dbg("Client disconnected")
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
	dbg("Connected to %s %s", c.Config.Name, c.Config.Path)
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

	methodPath := c.Config.Iface + "." + method

	dbg("Call %s( %v )", methodPath, args)

	return c.dbusObject.Call(methodPath, flags, args...)
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
	return c.dbusObject.Call("org.freedesktop.DBus.Properties.Set", 0, c.Config.Iface, p, v).Store()
}

//GetProperties load all the properties for an interface
func (c *Client) GetProperties(props interface{}) error {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}

	dbg("Loading properties for %s", c.Config.Iface)

	result := make(map[string]dbus.Variant)
	err := c.dbusObject.Call("org.freedesktop.DBus.Properties.GetAll", 0, c.Config.Iface).Store(&result)
	if err != nil {
		return err
	}

	return util.MapToStruct(props, result)
}

func getMatchString(path string, iface string) string {
	return "type='signal',interface='" + iface + "',path='" + path + "'"
}

//Register for signals
func (c *Client) Register(path string, iface string) (chan *dbus.Signal, error) {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return nil, err
		}
	}

	matchstr := getMatchString(path, iface)
	dbg("Match on %s", matchstr)
	c.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchstr)

	channel := make(chan *dbus.Signal, 1)
	c.conn.Signal(channel)

	return channel, nil
}

//Unregister for signals
func (c *Client) Unregister(path string, iface string) error {
	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}
	matchstr := getMatchString(path, iface)
	dbg("Match on %s", matchstr)
	c.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, matchstr)

	return nil
}
