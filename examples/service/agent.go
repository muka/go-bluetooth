package main

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	log "github.com/sirupsen/logrus"
)

func RegisterAgent(agent profile.Agent1Interface, caps string) error {
	//agent_path := AgentDefaultRegisterPath // we use the default path
	agent_path := agent.RegistrationPath() // we use the default path
	log.Infof("Agent path: %s", agent_path)

	// Register agent
	am := profile.NewAgentManager1(agent_path)

	// Export the Go interface to DBus
	err := am.ExportGoAgentToDBus(agent)
	if err != nil {
		return err
	}

	// Register the exported interface as application agent via AgenManager API
	err = am.RegisterAgent(agent_path, caps)
	if err != nil {
		return err
	}

	// Set the new application agent as Default Agent
	err = am.RequestDefaultAgent(agent_path)
	if err != nil {
		return err
	}

	return nil
}

func createAgent() (*Agent, error) {
	agent := new(Agent)
	agent.BusName = bluez.OrgBluezInterface
	agent.AgentInterface = bluez.Agent1Interface
	agent.AgentPath = agentObjectPath

	return agent, RegisterAgent(agent, profile.AGENT_CAP_KEYBOARD_DISPLAY)
}

func setTrusted(path string) {

	devices, err := api.GetDevices()
	if err != nil {
		log.Error(err)
		return
	}

	for _, v := range devices {
		if strings.Contains(v.Path, path) {
			log.Infof("Found device %s", path)
			dev1, _ := v.GetClient()
			err := dev1.SetProperty("Trusted", true)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

type Agent struct {
	BusName        string
	AgentInterface string
	AgentPath      string
}

func (self *Agent) Release() *dbus.Error {
	return nil
}

func (self *Agent) RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) {
	log.Info("RequestPinCode", device)
	setTrusted(string(device))
	return "0000", nil
}

func (self *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	log.Info(fmt.Sprintf("DisplayPinCode (%s, %s)", device, pincode))
	return nil
}

func (self *Agent) RequestPasskey(device dbus.ObjectPath) (passkey uint32, err *dbus.Error) {
	setTrusted(string(device))
	return 1024, nil
}

func (self *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	log.Info(fmt.Sprintf("DisplayPasskey %s, %06u entered %u", device, passkey, entered))
	return nil
}

func (self *Agent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	log.Info(fmt.Sprintf("RequestConfirmation (%s, %06d)", device, passkey))
	setTrusted(string(device))
	return nil
}

func (self *Agent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	log.Infof("RequestAuthorization (%s)\n", device)
	return nil
}

func (self *Agent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	log.Infof("AuthorizeService (%s, %s)", device, uuid) //directly authrized
	return nil
}

func (self *Agent) Cancel() *dbus.Error {
	log.Info("Cancel")
	return nil
}

func (self *Agent) RegistrationPath() string {
	return self.AgentPath
}

func (self *Agent) InterfacePath() string {
	return self.AgentInterface
}
