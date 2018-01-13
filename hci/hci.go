package hci

import (
	"github.com/paypal/gatt/linux"
)

//Do is a dev placeholder
func Do() error {
	devID := 0
	chk := false
	maxConn := 3

	hci, err := linux.NewHCI(devID, chk, maxConn)
	if err != nil {
		return err
	}
	defer hci.Close()

	err = hci.SetAdvertiseEnable(true)
	if err != nil {
		return err
	}

	err = hci.SetScanEnable(true, true)
	if err != nil {
		return err
	}

	return nil
}
