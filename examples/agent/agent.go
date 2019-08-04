package agent_example

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/linux/btmgmt"
	log "github.com/sirupsen/logrus"
)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)
func Run(deviceAddress, adapterID string) error {

	defer api.Exit()

	a := btmgmt.NewBtMgmt(adapterID)
	err := a.Reset()
	if err != nil {
		return err
	}

	ag := agent.NewSimpleAgent()
	err = agent.ExposeAgent(ag, agent.CapKeyboardDisplay, true)
	if err != nil {
		return fmt.Errorf("SimpleAgent: %s", err)
	}

	devices, err := api.GetDevices(adapterID)
	if err != nil {
		return fmt.Errorf("GetDevices: %s", err)
	}

	for _, dev := range devices {

		if dev.Properties.Address != deviceAddress {
			continue
		}

		// log.Info(i, v.Path)
		log.Infof("Pairing with %s", dev.Properties.Address)

		err := dev.Pair()
		if err != nil {
			return fmt.Errorf("Pair failed: %s", err)
		}

		log.Info("Pair succeed, connecting...")
		agent.SetTrusted(adapterID, dbus.ObjectPath(dev.Path))

		err = dev.Connect()
		if err != nil {
			return fmt.Errorf("Connect failed: %s", err)
		}

	}

	log.Info("Working...")
	select {}
}
