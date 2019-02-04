package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

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
