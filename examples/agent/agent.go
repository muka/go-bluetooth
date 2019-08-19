package agent_example

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	log "github.com/sirupsen/logrus"
)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)
func Run(deviceAddress, adapterID string) error {

	defer api.Exit()

	//Connect DBus System bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	ag := agent.NewSimpleAgent()
	err = agent.ExposeAgent(conn, ag, agent.CapKeyboardDisplay, true)
	if err != nil {
		return fmt.Errorf("SimpleAgent: %s", err)
	}

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	devices, err := a.GetDevices()
	if err != nil {
		return fmt.Errorf("GetDevices: %s", err)
	}

	found := false
	for _, dev := range devices {

		if dev.Properties.Address != deviceAddress {
			continue
		}

		if dev.Properties.Paired {
			continue
		}

		found = true
		// log.Info(i, v.Path)
		log.Infof("Pairing with %s", dev.Properties.Address)

		err := dev.Pair()
		if err != nil {
			return fmt.Errorf("Pair failed: %s", err)
		}

		log.Info("Pair succeed, connecting...")
		agent.SetTrusted(adapterID, dev.Path())

		err = dev.Connect()
		if err != nil {
			return fmt.Errorf("Connect failed: %s", err)
		}

	}

	if !found {
		return fmt.Errorf("No device found that need to be paired on %s", adapterID)
	}

	log.Info("Working...")
	select {}
}
