package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/linux/btmgmt"
	log "github.com/sirupsen/logrus"
)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)

const adapterID = "hci0"
const dev_mac = "88:B4:A6:6F:12:EF"

const logLevel = log.DebugLevel

const (
	BUS_NAME        = "org.bluez"
	AGENT_INTERFACE = "org.bluez.Agent1"
	AGENT_PATH      = "/gameshell/bleagentgo"
)

func main() {

	log.SetLevel(logLevel)

	defer api.Exit()

	a := btmgmt.NewBtMgmt(adapterID)
	err := a.Reset()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = InitAgent()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	devices, err := api.GetDevices()
	if err != nil {
		log.Error("GetDevices", err)
		return
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
				set_trusted(dev_mac_path)
				v.Connect()
			} else {
				log.Error("Pair failed:", err)
			}
		}
	}

	log.Info("Working...")
	select {}
}

func RegisterAgent(agent profile.Agent1Interface, caps string) error {
	//agent_path := AgentDefaultRegisterPath // we use the default path
	agent_path := agent.RegistrationPath() // we use the default path
	log.Info("The Agent Path: ", agent_path)

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

func InitAgent() error {
	agent := new(Agent)
	agent.BusName = BUS_NAME
	agent.AgentInterface = AGENT_INTERFACE
	agent.AgentPath = AGENT_PATH

	return RegisterAgent(agent, profile.AGENT_CAP_KEYBOARD_DISPLAY)
}

func set_trusted(path string) {
	devices, err := api.GetDevices()
	if err != nil {
		log.Error(err)
		return
	}

	for i, v := range devices {
		log.Info(i, v.Path)
		if strings.Contains(v.Path, path) {
			log.Info("Found device")
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
	set_trusted(string(device))
	return "0000", nil
}

func (self *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	log.Info(fmt.Sprintf("DisplayPinCode (%s, %s)", device, pincode))
	return nil
}

func (self *Agent) RequestPasskey(device dbus.ObjectPath) (passkey uint32, err *dbus.Error) {
	set_trusted(string(device))
	return 1024, nil
}

func (self *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	log.Info(fmt.Sprintf("DisplayPasskey %s, %06u entered %u", device, passkey, entered))
	return nil
}

func (self *Agent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	log.Info(fmt.Sprintf("RequestConfirmation (%s, %06d)", device, passkey))
	set_trusted(string(device))
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
