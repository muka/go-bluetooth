package agent_example

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/linux/btmgmt"
	log "github.com/sirupsen/logrus"
)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)

const (
	BUS_NAME   = "org.bluez"
	AGENT_PATH = "/gameshell/bleagentgo"
)

func Run(dev_mac, adapterID string) error {

	defer api.Exit()

	a := btmgmt.NewBtMgmt(adapterID)
	err := a.Reset()
	if err != nil {
		return err
	}

	err = InitAgent()
	if err != nil {
		return fmt.Errorf("InitAgent: %s", err)
	}

	devices, err := api.GetDevices()
	if err != nil {
		return fmt.Errorf("GetDevices: %s", err)
	}

	log.Info(devices)

	for i, v := range devices {
		dev_mac_path := strings.Replace(dev_mac, ":", "_", -1)
		if strings.Contains(v.Path, dev_mac_path) {
			log.Info(i, v.Path)
			log.Info("Pairing...")
			err := v.Pair()
			if err == nil {
				log.Info("Pair succeed,Connecting...")
				setTrusted(dev_mac_path)
				v.Connect()
			} else {
				return fmt.Errorf("Pair failed: %s", err)
			}
		}
	}

	log.Info("Working...")
	select {}
}

func RegisterAgent(ag agent.Agent1Client, caps string) error {

	agentPath := ag.RegistrationPath() // we use the default path
	log.Info("Agent Path: ", agentPath)

	// Register agent
	am, err := agent.NewAgentManager1()
	if err != nil {
		return err
	}

	// Export the Go interface to DBus
	err = agent.ExportAgent(ag)
	if err != nil {
		return err
	}

	// Register the exported interface as application agent via AgenManager API
	err = am.RegisterAgent(dbus.ObjectPath(agentPath), caps)
	if err != nil {
		return err
	}

	// Set the new application agent as Default Agent
	err = am.RequestDefaultAgent(dbus.ObjectPath(agentPath))
	if err != nil {
		return err
	}

	return nil
}

func InitAgent() error {

	ag := new(Agent)
	ag.BusName = BUS_NAME
	ag.AgentInterface = agent.Agent1Interface
	ag.AgentPath = AGENT_PATH

	return RegisterAgent(ag, agent.AGENT_CAP_KEYBOARD_DISPLAY)
}

func setTrusted(path string) {

	devices, err := api.GetDevices()
	if err != nil {
		log.Error(err)
		return
	}

	for i, v := range devices {
		log.Info(i, v.Path)
		if strings.Contains(v.Path, path) {
			log.Info("Found device")
			dev1, err := v.GetClient()
			if err != nil {
				log.Error(err)
				return
			}
			err = dev1.SetProperty("Trusted", true)
			if err != nil {
				log.Error(err)
				break
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
	log.Info(fmt.Sprintf("DisplayPasskey %s, %06d entered %d", device, passkey, entered))
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
