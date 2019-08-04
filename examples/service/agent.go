package service_example

import (
	"github.com/muka/go-bluetooth/bluez/profile/agent"
)

func createAgent() (agent.Agent1Client, error) {
	a := agent.NewSimpleAgent()
	return a, agent.ExposeAgent(a, agent.CapKeyboardDisplay, true)
}
