package agent

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

//ExportAgent exports the xml of a go agent to dbus
func ExportAgent(agentInstance Agent1Interface) error {

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

	// Export Introspectable for the given agent instance
	err = conn.Export(introspect.NewIntrospectable(node), dbus.ObjectPath(targetPath), bluez.Introspectable)
	if err != nil {
		return err
	}

	return nil
}
