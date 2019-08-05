package bluez

// PropertyChanged indicates that a change is notified
type PropertyChanged struct {
	Interface string
	Name      string
	Value     interface{}
}
