package service

import "github.com/muka/go-bluetooth/bluez/profile/agent"

func (app *App) createAgent() (agent.Agent1Client, error) {
	a := agent.NewDefaultSimpleAgent()
	return a, nil
}

// Expose app agent on DBus
func (app *App) ExposeAgent(caps string, setAsDefaultAgent bool) error {
	return agent.ExposeAgent(app.DBusConn(), app.agent, caps, setAsDefaultAgent)
}
