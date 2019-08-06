package agent

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	log "github.com/sirupsen/logrus"
)

//All agent capabilities
const (
	CapDisplayOnly     = "DisplayOnly"
	CapDisplayYesNo    = "DisplayYesNo"
	CapKeyboardOnly    = "KeyboardOnly"
	CapNoInputNoOutput = "NoInputNoOutput"
	CapKeyboardDisplay = "KeyboardDisplay"
)

type Agent1Client interface {
	Release() *dbus.Error                                                    // Callback doesn't trigger on unregister
	RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) // Triggers for pairing when SSP is off and cap != CAP_NO_INPUT_NO_OUTPUT
	DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error
	RequestPasskey(device dbus.ObjectPath) (passkey uint32, err *dbus.Error) // SSP on, toolz.AGENT_CAP_KEYBOARD_ONLY
	DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error
	RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error
	RequestAuthorization(device dbus.ObjectPath) *dbus.Error
	AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error
	Cancel() *dbus.Error
	Path() dbus.ObjectPath
	Interface() string
}

// SetTrusted lookup for a device by object path and set it to trusted
func SetTrusted(adapterID string, devicePath dbus.ObjectPath) error {

	log.Debugf("Trust device %s on %s", devicePath, adapterID)

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	devices, err := a.GetDevices()
	if err != nil {
		return err
	}

	path := string(devicePath)
	for _, dev := range devices {
		if strings.Contains(string(dev.Path()), path) {
			log.Debugf("Trust device at %s", path)
			err := dev.SetProperty("Trusted", true)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Cannot trust device %s, not found", path)
}

// ExposeAgent expose an Agent1 implementation to DBus and set as default agent
func ExposeAgent(ag Agent1Client, caps string, setAsDefaultAgent bool) error {

	// Register agent
	am, err := NewAgentManager1()
	if err != nil {
		return err
	}

	// Export the Go interface to DBus
	err = exportAgent(ag)
	if err != nil {
		return err
	}

	// Register the exported interface as application agent via AgenManager API
	err = am.RegisterAgent(ag.Path(), caps)
	if err != nil {
		return err
	}

	if setAsDefaultAgent {
		// Set the new application agent as Default Agent
		err = am.RequestDefaultAgent(ag.Path())
		if err != nil {
			return err
		}
	}

	return nil
}

//ExportAgent exports the xml of a go agent to dbus
func exportAgent(agentInstance Agent1Client) error {

	//Connect DBus System bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	targetPath := agentInstance.Path()
	agentInterfacePath := agentInstance.Interface()

	//Export the given agent to the given path as interface "org.bluez.Agent1"
	err = conn.Export(agentInstance, targetPath, agentInterfacePath)
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
	err = conn.Export(introspect.NewIntrospectable(node), targetPath, bluez.Introspectable)
	if err != nil {
		return err
	}

	return nil
}
