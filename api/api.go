// package api wraps an high level API to simplify interaction
package api

import (
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

var adapters = map[string]*adapter.Adapter1{}

//Exit performs a clean exit
func Exit() error {

	for _, a := range adapters {
		a.Close()
	}

	adapters = map[string]*adapter.Adapter1{}

	return bluez.CloseConnections()
}

func GetAdapter(adapterID string) (*adapter.Adapter1, error) {

	if _, ok := adapters[adapterID]; ok {
		return adapters[adapterID], nil
	}

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return nil, err
	}

	adapters[adapterID] = a

	return a, nil
}

func GetDefaultAdapter() (*adapter.Adapter1, error) {
	return GetAdapter(GetDefaultAdapterID())
}

func GetDefaultAdapterID() string {
	return adapter.GetDefaultAdapterID()
}
