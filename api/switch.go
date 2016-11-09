package api

import (
	"errors"

	"github.com/muka/bluez-client/linux"
)

var rfkill = linux.NewRFKill()

// ToggleDevice Swap Off/On a device
func ToggleDevice(identifier string) error {
	err := TurnOffDevice(identifier)
	if err != nil {
		return err
	}
	return TurnOnDevice(identifier)
}

// TurnOnDevice Enable a rfkill managed device
func TurnOnDevice(identifier string) error {
	if !rfkill.IsInstalled() {
		return errors.New("rfkill is not available")
	}
	if rfkill.IsSoftBlocked(identifier) {
		err := rfkill.SoftUnblock(identifier)
		if err != nil {
			return err
		}
	}
	if rfkill.IsHardBlocked(identifier) {
		return errors.New("Device is hard locked, check for a physical switch to enable it")
	}
	return nil
}

// TurnOffDevice Enable a rfkill managed device
func TurnOffDevice(identifier string) error {
	if !rfkill.IsInstalled() {
		return errors.New("rfkill is not available")
	}
	if !rfkill.IsSoftBlocked(identifier) {
		err := rfkill.SoftBlock(identifier)
		if err != nil {
			return err
		}
	}
	return nil
}
