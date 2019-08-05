package api

import (
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/linux/btmgmt"
)

//Exit performs a clean exit
func Exit() {
	bluez.CloseConnections()
}

func ResetController(adapterID string) error {
	a := btmgmt.NewBtMgmt(adapterID)
	err := a.Reset()
	if err != nil {
		return err
	}
	return nil
}

func GetAdapter(adapterID string) (*adapter.Adapter1, error) {
	return adapter.GetAdapter(adapterID)
}

func GetDefaultAdapter() (*adapter.Adapter1, error) {
	return adapter.GetDefaultAdapter()
}
