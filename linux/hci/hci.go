package hci

import (
	"github.com/currantlabs/ble/linux/hci"
)

//Do is a dev placeholder
func Do() error {

	h, err := hci.NewHCI()
	if err != nil {
		return err
	}

	err = h.Init()
	if err != nil {
		return err
	}

	err = h.Scan(false)
	if err != nil {
		return err
	}

	return nil
}
