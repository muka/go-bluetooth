package adapter

// GetAdapterID return the Id of the adapter
func (a *Adapter1) GetAdapterID() (string, error) {
	return ParseAdapterID(a.Path())
}
