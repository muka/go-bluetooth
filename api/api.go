// package api wraps an high level API to simplify interaction
package api

import (
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

//Exit performs a clean exit
func Exit() error {
	return bluez.CloseConnections()
}

func GetAdapter(adapterID string) (*adapter.Adapter1, error) {
	return adapter.GetAdapter(adapterID)
}

func GetDefaultAdapter() (*adapter.Adapter1, error) {
	return adapter.GetDefaultAdapter()
}

func GetDefaultAdapterID() string {
	return adapter.GetDefaultAdapterID()
}
