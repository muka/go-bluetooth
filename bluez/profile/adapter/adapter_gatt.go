package adapter

import (
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//GetGattManager return a GattManager1 instance
func (a *Adapter1) GetGattManager() (*gatt.GattManager1, error) {
	adapterID, err := ParseAdapterID(a.Path())
	if err != nil {
		return nil, err
	}
	return gatt.NewGattManager1FromAdapterID(adapterID)
}
