package bluez

import (
	"log"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/util"
)

const (
	// SessionBus uses the session bus
	SessionBus = 0
	// SystemBus uses the system bus
	SystemBus = 1
)

var conns = make([]*dbus.Conn, 2)

// Config pass configuration to a DBUS client
type Config struct {
	Name  string
	Iface string
	Path  string
	Bus   int
}

// NewClient create a new client
func NewClient(config *Config) *Client {

	c := new(Client)

	c.Config = config
	c.logger = util.NewLogger("client")

	return c
}

// Client implement a DBus client
type Client struct {
	logger     *log.Logger
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
		c.logger.Println("Client disconnected")
	}
}

// Connect connect to DBus
func (c *Client) Connect() error {

	var getConn = func() (*dbus.Conn, error) {
		switch c.Config.Bus {
		case SystemBus:
			{
				if conns[SystemBus] == nil {
					c.logger.Println("Connecting to SystemBus")
					conn, err := dbus.SystemBus()
					if err != nil {
						return nil, err
					}
					conns[SystemBus] = conn
				}
				return conns[SystemBus], nil
			}
		case SessionBus:
			{
				if conns[SessionBus] == nil {
					c.logger.Println("Connecting to SessionBus")
					conn, err := dbus.SessionBus()
					if err != nil {
						return nil, err
					}
					conns[SessionBus] = conn
				}
				return conns[SessionBus], nil
			}
		default:
			{
				c.logger.Println("TODO: Unknown Bus, handle other types!")
				return nil, nil
			}
		}
	}

	dbusConn, err := getConn()
	if err != nil {
		return err
	}

	c.conn = dbusConn
	c.dbusObject = c.conn.Object(c.Config.Name, dbus.ObjectPath(c.Config.Path))

	c.logger.Printf("Connected to %s %s\n", c.Config.Name, c.Config.Path)
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

	callArgs := args
	c.logger.Printf("Call %s( %v )\n", methodPath, callArgs)

	return c.dbusObject.Call(methodPath, flags, callArgs...)
}

//GetProperty return a property value
func (c *Client) GetProperty(p string) (dbus.Variant, error) {
	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return dbus.Variant{}, err
		}
	}
	return c.dbusObject.GetProperty(p)
}

//GetProperties load all the properties for an interface
func (c *Client) GetProperties(props interface{}) error {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}

	c.logger.Printf("Loading properties for %s", c.Config.Iface)

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
	c.logger.Printf("Match on %s", matchstr)
	c.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchstr)

	channel := make(chan *dbus.Signal, 100)
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
	c.logger.Printf("Match on %s", matchstr)
	c.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, matchstr)
	return nil
}
