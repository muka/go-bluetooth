package api

import (
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/hw"
)

//Exit performs a clean exit
func Exit() error {
	return bluez.CloseConnections()
}

func ResetController(adapterID string) error {
	return hw.Reset(adapterID)
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
