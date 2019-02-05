package profile

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
)

//All agent capabilities
const (
	AGENT_CAP_DISPLAY_ONLY       = "DisplayOnly"
	AGENT_CAP_DISPLAY_YES_NO     = "DisplayYesNo"
	AGENT_CAP_KEYBOARD_ONLY      = "KeyboardOnly"
	AGENT_CAP_NO_INPUT_NO_OUTPUT = "NoInputNoOutput"
	AGENT_CAP_KEYBOARD_DISPLAY   = "KeyboardDisplay"
)

type Agent1Interface interface {
	Release() *dbus.Error                                                    // Callback doesn't trigger on unregister
	RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) // Triggers for pairing when SSP is off and cap != CAP_NO_INPUT_NO_OUTPUT
	DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error
	RequestPasskey(device dbus.ObjectPath) (passkey uint32, err *dbus.Error) // SSP on, toolz.AGENT_CAP_KEYBOARD_ONLY
	DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error
	RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error
	RequestAuthorization(device dbus.ObjectPath) *dbus.Error
	AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error
	Cancel() *dbus.Error
	RegistrationPath() string
	InterfacePath() string
}

// NewAgentManager1 create a new AgentManager1 client
func NewAgentManager1(hostID string) *AgentManager1 {
	a := new(AgentManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.AgentManager1",
			Path:  "/org/bluez",
			Bus:   bluez.SystemBus,
		},
	)
	return a
}

// AgentManager1 client
type AgentManager1 struct {
	client *bluez.Client
}

// Close the connection
func (a *AgentManager1) Close() {
	a.client.Disconnect()
}

//RegisterAgent registers an agent handler
func (a *AgentManager1) RegisterAgent(agent string, capability string) error {
	return a.client.Call("RegisterAgent", 0, dbus.ObjectPath(agent), capability).Store()
}

//RequestDefaultAgent requests to make the application agent the default agent
func (a *AgentManager1) RequestDefaultAgent(agent string) error {
	return a.client.Call("RequestDefaultAgent", 0, dbus.ObjectPath(agent)).Store()
}

//UnregisterAgent unregisters the agent that has been previously registered
func (a *AgentManager1) UnregisterAgent(agent string) error {
	return a.client.Call("UnregisterAgent", 0, dbus.ObjectPath(agent)).Store()
}

//ExportGoAgentToDBus exports the xml of go agent to dbus
func (a *AgentManager1) ExportGoAgentToDBus(agentInstance Agent1Interface) error {

	//Connect DBus System bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	targetPath := agentInstance.RegistrationPath()
	agentInterfacePath := agentInstance.InterfacePath()

	//Export the given agent to the given path as interface "org.bluez.Agent1"
	err = conn.Export(agentInstance, dbus.ObjectPath(targetPath), agentInterfacePath)
	if err != nil {
		return err
	}

	// Create  Introspectable for the given agent instance
	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			// Introspect
			introspect.IntrospectData,
			// Properties
			prop.IntrospectData,
			// org.bluez.Agent1
			{
				Name:    agentInterfacePath,
				Methods: introspect.Methods(agentInstance),
			},
		},
	}
	//fmt.Println(node)

	// Export Introspectable for the given agent instance
	err = conn.Export(introspect.NewIntrospectable(node), dbus.ObjectPath(targetPath), "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	return nil
}
