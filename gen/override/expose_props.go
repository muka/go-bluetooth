package override

var ExposePropertiesInterface = map[string]bool{
	"org.bluez.AgentManager1":   false,
	"org.bluez.Agent1":          false,
	"org.bluez.ProfileManager1": false,
	"org.bluez.Profile1":        false,
}

// ExposeProperties expose Properties interface to the struct
func ExposeProperties(iface string) bool {
	if val, ok := ExposePropertiesInterface[iface]; ok {
		return val
	}
	return true
}
