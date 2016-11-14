package linux

import (
	"errors"
	sddbus "github.com/coreos/go-systemd/dbus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

type result string

const (
	done       result = "done"
	canceled   result = "canceled"
	timeout    result = "timeout"
	failed     result = "failed"
	dependency result = "dependency"
	skipped    result = "skipped"
)

var systemdConn *sddbus.Conn
var ch = make(chan string)
var inprogress = false
var callback func(err error)

const bluetoothUnitName = "bluetooth.service"

// @see https://github.com/coreos/go-systemd/blob/master/dbus/methods.go#L65
const defaultMode = "replace"

func getConnection() (*sddbus.Conn, error) {
	return sddbus.NewConnection(func() (*dbus.Conn, error) {
		return bluez.GetConnection(bluez.SystemBus)
	})
}

func watchOperation() {
	inprogress = true
	for response := range ch {
		switch result(response) {
		case done:
			// success
			callback(nil)
			break
		case canceled:
		case timeout:
		case failed:
		case dependency:
		case skipped:
			//fail
			callback(errors.New("Unit operation failed: " + response))
			break
		}
		inprogress = false
		callback = nil
		// free resources
		systemdConn.Close()
		systemdConn = nil
		return
	}
}

//StartBluetoothUnit by its systemd unit
func StartBluetoothUnit(fn func(err error)) (int, error) {
	systemdConn, err := getConnection()
	if err != nil {
		return 0, err
	}
	if inprogress {
		return 0, nil
	}
	callback = fn
	go watchOperation()
	systemdConn.StartUnit(bluetoothUnitName, defaultMode, ch)
	return 0, nil
}

//StopBluetoothUnit by its systemd unit
func StopBluetoothUnit(fn func(err error)) (int, error) {
	systemdConn, err := getConnection()
	if err != nil {
		return 0, err
	}
	if inprogress {
		return 0, nil
	}
	callback = fn
	go watchOperation()
	return systemdConn.StopUnit(bluetoothUnitName, defaultMode, ch)
}

//RestartBluetoothUnit by its systemd unit
func RestartBluetoothUnit(fn func(err error)) (int, error) {
	systemdConn, err := getConnection()
	if err != nil {
		return 0, err
	}
	if inprogress {
		return 0, nil
	}
	callback = fn
	go watchOperation()
	return systemdConn.RestartUnit(bluetoothUnitName, defaultMode, ch)
}
